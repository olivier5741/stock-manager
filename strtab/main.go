package strtab

import (
	"fmt"
)

var (
	rowSizeErr = fmt.Errorf("The size of the row doesn't match the size of content row")
)

type T struct {
	RowHeader []string
	ColHeader []string
	Content   [][]string
}

func NewT() *T {
	t := T{make([]string, 0), make([]string, 0), make([][]string, 0)}
	return &t
}

func (t *T) AddRowHeader(r []string) error {
	if len(t.Content) != 0 && len(t.Content[0]) != len(r) {
		return rowSizeErr
	}
	copy(t.RowHeader, r)
	return nil
}

func (t *T) AddColHeader(c []string) error {
	if len(t.Content) != 0 && len(t.Content) != len(c) {
		return rowSizeErr
	}
	copy(t.ColHeader, c)
	return nil
}

func (t *T) addRow(v []string) error {
	if len(t.Content) != 0 && len(t.Content[0]) != len(v) {
		return rowSizeErr
	}
	if len(t.Content) == 0 {
		t.ColHeader = make([]string, len(v))
	}
	t.Content = append(t.Content, v)
	return nil
}

func (t *T) AddRowWithoutHeader(v []string) error {
	return t.AddRowWithHeader("", v)
}

func (t *T) AddRowWithHeader(h string, v []string) error {
	err := t.addRow(v)
	if err != nil {
		return err
	}
	t.RowHeader = append(t.RowHeader, h)
	return nil
}

func (t *T) addCol(v []string) error {
	if len(t.Content) != 0 && len(t.Content) != len(v) {
		return fmt.Errorf("The size of the column doesn't match the size of content column")
	}
	if len(t.Content) == 0 {
		t.RowHeader = make([]string, len(v))
		for i := 0; i < len(v); i++ {
			t.Content = append(t.Content, make([]string, 0))
		}
	}

	for k, i := range t.Content {
		t.Content[k] = append(i, v[k])
	}

	return nil
}

func (t *T) AddColWithoutHeader(v []string) error {
	return t.AddColWithHeader("", v)
}

func (t *T) AddColWithHeader(h string, v []string) error {
	err := t.addCol(v)
	if err != nil {
		return err
	}
	t.ColHeader = append(t.ColHeader, h)
	return nil
}

func (t T) GetContent() [][]string {
	out := make([][]string, 0)
	for _, r := range t.Content {
		newRow := make([]string, len(r))
		copy(newRow, r)
		out = append(out, newRow)
	}
	return out
}

func (t T) GetContentWithColHeader() [][]string {
	out := make([][]string, 0)
	out = append(out, t.ColHeader)
	for _, c := range t.GetContent() {
		out = append(out, c)
	}
	return out
}

func (t T) GetContentWithRowHeader() [][]string {
	out := make([][]string, 0)
	for i, c := range t.GetContent() {
		out = append(out, append([]string{t.RowHeader[i]}, c...))
	}
	return out
}
