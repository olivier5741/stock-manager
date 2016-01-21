package strtab

import (
	"testing"
)

func Test(t *testing.T) {
	tab := NewT()
	tab.AddRowWithHeader("iso", []string{"22", "12"})
	tab.AddRowWithHeader("aspi", []string{"2", "3"})
	tab.AddRowWithHeader("busco", []string{"1", "11"})
	t.Log(tab.GetContentWithColHeader())
	t.Log(tab.GetContentWithRowHeader())

	tab1 := NewT()
	tab1.AddColWithHeader("iso", []string{"22", "12"})
	tab1.AddColWithHeader("aspi", []string{"2", "3"})
	tab1.AddColWithHeader("busco", []string{"1", "11"})
	t.Log(tab1.GetContentWithColHeader())
	t.Log(tab1.GetContentWithRowHeader())
}
