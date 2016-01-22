package strtab

import (
	"testing"
)

func Test(t *testing.T) {
	tab := NewT()

	tab.AddColHeader([]string{"min", "max"})
	tab.AddWithHeader([]string{"iso", "2", "3"},
		[]string{"aspi", "22", "12"},
		[]string{"busco", "1", "11"})
	t.Log(tab.GetContentWithRowHeader())
	t.Log(tab.GetContentWithColHeader())
	t.Log(tab.GetContentWithHeaders(true))

	tab.Transpose()
	t.Log(tab.GetContentWithRowHeader())
	t.Log(tab.GetContentWithColHeader())
	t.Log(tab.GetContentWithHeaders(true))
	t.Log(tab.String())

	// TODO : put real test
}
