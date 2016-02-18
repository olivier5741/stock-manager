package sheet

import(
	"testing"
	"github.com/olivier5741/stock-manager/port/sheet/osfile"
)

func TestAllSheets(t *testing.T) {
	d := osfile.OsFile{"../../cmd/main/"}
	s := AllSheets(d)
	t.Log(s)
	t.Log(NewSheet(s[0],d))
}