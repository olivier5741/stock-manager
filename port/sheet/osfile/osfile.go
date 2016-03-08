package osfile

import(
	"io/ioutil"
	"io"
	"os"
	"strings"
	log "github.com/Sirupsen/logrus"
	"github.com/olivier5741/stock-manager/asset"
	"github.com/olivier5741/stock-manager/port/sheet"
)

var(
	extension = ".csv"
	Tr = func(s string) string {return s}
)

type OsFile struct{
	Dir string // "./" as default
}

func (o OsFile) GetAll() []sheet.BasicFilename{

	var out []sheet.BasicFilename
	files, err := ioutil.ReadDir(o.Dir)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"dir": o.Dir,
		}).Error(asset.Tr("error_dir_read"))
	}

	for _, file := range files {

		if !strings.HasSuffix(file.Name(), extension) {
			continue
		}

		id := strings.TrimSuffix(file.Name(),extension)
		out = append(out,sheet.BasicFilename{id,id})	
	}
	return out
}

func (o OsFile) NewReader(name sheet.BasicFilename) io.ReadCloser {

	r, err := os.OpenFile(o.Dir + name.Name + extension, os.O_CREATE, 0666)
	
	if err != nil {
		log.WithFields(log.Fields{
		 	"filepath": o.Dir + name.Name + extension,
		 	"err":      err,
		}).Error(asset.Tr("error_file_open"))
	}

	return r
}

func (o OsFile) NewWriter(name sheet.BasicFilename) io.WriteCloser {
	w, err := os.Create(o.Dir + name.Name + extension)

	if err != nil {
		log.WithFields(log.Fields{
			"filepath": o.Dir + name.Name + extension,
			"err":      err,
		}).Error(asset.Tr("create_file_error"))
	}

	return w

}
