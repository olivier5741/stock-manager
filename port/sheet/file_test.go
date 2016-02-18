package sheet

import(
"testing"
)

func TestNewFilename(t *testing.T) {
	f1 := "2016-01-28-n1-entree"
	out1, err1 := NewFilename(f1)
	t.Log(out1)
	t.Log(err1)

	f2 := "en attente-2016-02-18-n3-inventaire"
	out2, err2 := NewFilename(f2)
	t.Log(out2)
	t.Log(err2)
}