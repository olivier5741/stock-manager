package osfile

import(
	"io/ioutil"
	"io"
	"os"
	"strings"
)

var(
	extension = ".csv"
)

type OsFile struct{
	Dir string // "./" as default
}

func (o OsFile) GetAll() []string {

	var out []string
	files, err1 := ioutil.ReadDir(o.Dir)

	if err1 != nil {
		// log.WithFields(log.Fields{
		// 	"err": err1,
		// }).Error(Tr("error_dir_read"))
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), extension) {
			continue
		}
		out = append(out,strings.TrimSuffix(file.Name(),extension))
	}

	return out
}

func (o OsFile) NewReader(name string) io.ReadCloser {

	r, err := os.OpenFile(o.Dir + name + extension, os.O_CREATE, 0666)
	
	if err != nil {
		// log.WithFields(log.Fields{
		// 	"filename": filename,
		// 	"err":      err,
		// }).Error(Tr("error_file_open"))
	}

	return r
}

func (o OsFile) NewWriter(name string) io.WriteCloser {
	w, err := os.Create(o.Dir + name + extension)

	if err != nil {
		// log.WithFields(log.Fields{
		// 	"path":     t.Path,
		// 	"filename": t.Title,
		// 	"err":      err,
		// }).Error("create_file_error")
	}

	return w

}
