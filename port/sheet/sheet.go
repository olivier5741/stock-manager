package sheet

import(
	"github.com/olivier5741/strtab"
	"github.com/olivier5741/stock-manager/asset"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
	"encoding/csv"
	"fmt"
)

var (
	generatedPrefix = "g" + sep
)

type Filer interface {
	GetAll() []BasicFilename
	NewReader(f BasicFilename) io.ReadCloser
	NewWriter(f BasicFilename) io.WriteCloser
}

type Sheet struct {
	Name Filename
	Table  *strtab.T
	Render func(*strtab.T) [][]string
}

func AllSheets(o Filer) []Filename {
	var fs []Filename
	list := o.GetAll()
	for _,s := range list {
		if strings.HasPrefix(s.Name, generatedPrefix) ||
			strings.HasPrefix(s.Name, asset.Tr("file_prefix_draft")) {
			continue
		}

		name, err1 := NewFilename(s)
		if err1 != nil {
			log.WithFields(log.Fields{
				"filename": name.String(),
				"err":      err1,
			}).Error(asset.Tr("error_file_csv_unmarshal"))
		}

		fs = append(fs, name)
	}
	return fs
}

func NewSheet(f Filename, o Filer) Sheet {

	fmt.Println(f.Basic())

	r := o.NewReader(f.Basic())
	defer r.Close()

	rcsv := csv.NewReader(r)
	out, err := rcsv.ReadAll()

	if err != nil {
		log.WithFields(log.Fields{
			"filename": f.String(),
			"view":     out,
			"err":      err,
		}).Error(asset.Tr("error_file_csv_unmarshal"))
		return Sheet{f,nil,nil}
	}

	tab := strtab.NewT(out[0],out[1:]...)

	fmt.Println(tab)

	return Sheet{f,tab,nil}
}

func (s Sheet) Put(o Filer) {

	w := o.NewWriter(s.Name.Basic())
	defer w.Close()

	wcsv := csv.NewWriter(w)
 	wcsv.WriteAll(s.Render(s.Table))

 	if wcsv.Error() != nil {
		log.WithFields(log.Fields{
			"filename": s.Name.String(),
			"view":     s.Table,
			"err":      wcsv.Error(),
		}).Error(asset.Tr("error_file_csv_view_to"))
	}
}
