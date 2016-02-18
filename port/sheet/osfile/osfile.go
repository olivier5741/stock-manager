package osfile

import(
	"io/ioutil"
	"io"
	"os"
	"strings"
	log "github.com/Sirupsen/logrus"
)

var(
	extension = ".csv"
	Tr = func(s string) string {return s}
)

type OsFile struct{
	Dir string // "./" as default
}

func (o OsFile) GetAll() []string {

	var out []string
	files, err := ioutil.ReadDir(o.Dir)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"dir": o.Dir,
		}).Error(Tr("error_dir_read"))
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), extension) {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".go") {
			continue
		}
		out = append(out,strings.TrimSuffix(file.Name(),extension))
	}

	return out
}

func (o OsFile) NewReader(name string) io.ReadCloser {

	r, err := os.OpenFile(o.Dir + name + extension, os.O_CREATE, 0666)
	
	if err != nil {
		log.WithFields(log.Fields{
		 	"filepath": o.Dir + name + extension,
		 	"err":      err,
		}).Error(Tr("error_file_open"))
	}

	return r
}

func (o OsFile) NewWriter(name string) io.WriteCloser {
	w, err := os.Create(o.Dir + name + extension)

	if err != nil {
		log.WithFields(log.Fields{
			"filepath": o.Dir + name + extension,
			"err":      err,
		}).Error(Tr("create_file_error"))
	}

	return w

}
