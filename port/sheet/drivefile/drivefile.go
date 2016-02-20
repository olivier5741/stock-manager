package drivefile

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"io"
	"bytes"
	"golang.org/x/oauth2/google"

	"io/ioutil"
	"github.com/Sirupsen/logrus"
	"github.com/olivier5741/stock-manager/asset"
	"github.com/olivier5741/stock-manager/port/sheet"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("drive-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetService() *drive.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	return srv
}

type DriveFile struct {
	Parent string // 0BzIZ3dfuz-CERzU1c21XQks2cEU
	Service *drive.Service
	Names map[string]string
}

func (d DriveFile) GetAll() []sheet.BasicFilename {
	r, err := d.Service.Files.List().
		Q("'" + d.Parent + "' in parents and mimeType='application/vnd.google-apps.spreadsheet'").
		Fields("nextPageToken, files(id, name)").Do()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"dir": d.Parent,
		}).Error(asset.Tr("error_dir_read"))
	}

	var out []sheet.BasicFilename
	for _,e := range r.Files {
		out = append(out,sheet.BasicFilename{e.Name,e.Id})
		d.Names[e.Name] = e.Id
	}

	return out
}

func (d DriveFile) NewReader(h sheet.BasicFilename) io.ReadCloser {

	url := "https://docs.google.com/spreadsheets/d/" + h.UID + "/export?format=csv"

	response, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
		 	"filepath": h.UID,
		 	"err":      err,
		}).Error(asset.Tr("error_file_open"))
	}

	return response.Body
}

func (d DriveFile) NewWriter(h sheet.BasicFilename) io.WriteCloser {

	df := &drive.File{
		Name: h.Name,
		MimeType: "application/vnd.google-apps.spreadsheet",
		Parents: []string{d.Parent},
	}

	id, exist := d.Names[h.Name]

	return buff{&bytes.Buffer{}, GetService(), df, id, exist}

}

func (d DriveFile) Exists(name string) {

}

type buff struct {
	*bytes.Buffer
	Service *drive.Service
	File *drive.File
	id string
	exist bool
}

func (b buff) Close() error {

	if !b.exist{
		_, err := b.Service.Files.Create(b.File).Media(b.Buffer, googleapi.ContentType("text/csv")).Do()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"filepath": b.File.Name,
				"err":      err,
			}).Error(asset.Tr("create_file_error"))
		}
	}else{
		_, err := b.Service.Files.Update(b.id,&drive.File{}).Media(b.Buffer, googleapi.ContentType("text/csv")).Do()

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"filepath": b.File.Name,
				"err":      err,
			}).Error(asset.Tr("create_file_error")) // Should be update
		}
	}

	return nil
}