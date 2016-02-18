package sheet

import(
	"github.com/olivier5741/strtab"
	"github.com/olivier5741/stock-manager/port/sheet/osfile"
	"github.com/olivier5741/stock-manager/tr"
	"strings"
	"encoding/csv"
)

var (
	generatedPrefix = "g" + sep
)

type Sheet struct {
	Name Filename
	Table  *strtab.T
	Render func(*strtab.T) [][]string
}

func AllSheets(o osfile.OsFile) []string {
	var fs []string
	list := o.GetAll()
	for _,s := range list {
		// move this somewhere else
		if strings.HasPrefix(s, generatedPrefix) ||
			strings.HasPrefix(s, tr.Tr("file_prefix_draft")) {
			continue
		}

		fs = append(fs, s)
	}
	return fs
}

func NewSheet(s string, o osfile.OsFile) Sheet {
	name, err1 := NewFilename(s)
	if err1 != nil {
		//
	}

	r := o.NewReader(name.String())
	defer r.Close()

	rcsv := csv.NewReader(r)
	out, err2 := rcsv.ReadAll()

	if err2 != nil {
		//
	}

	tab := strtab.NewT(out[0],out[1:]...)

	return Sheet{name,tab,nil}
}

func (s Sheet) Put(o osfile.OsFile) {

	w := o.NewWriter(s.Name.String())
	defer w.Close()

	wcsv := csv.NewWriter(w)
 	wcsv.WriteAll(s.Render(s.Table))

 	if wcsv.Error() != nil {
		// log.WithFields(log.Fields{
		// 	"path":     t.Path,
		// 	"filename": t.Title,
		// 	"view":     t.Table,
		// 	"err":      w.Error(),
		// }).Error(Tr("error_file_csv_view_to"))
	}
}
