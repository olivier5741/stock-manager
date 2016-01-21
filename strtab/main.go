package strtab

import (
	"fmt"
)

var (
	rowSizeErr = fmt.Errorf("The size of the row doesn't match the size of content row")
	colSizeErr = fmt.Errorf("The size of the col doesn't match the size of content col")
)

type T struct {
	rowHeader []string
	colHeader []string
	content   [][]string
}

func (t T) String() (s string) {
	for _, l := range t.GetContentWithHeaders() {
		for _, v := range l {
			s += v + ","
		}
		s += "\n"
	}
	return
}

func NewT() *T {
	t := T{make([]string, 0), make([]string, 0), make([][]string, 0)}
	return &t
}

func (t *T) Transpose() {

	t.rowHeader, t.colHeader = t.colHeader, t.rowHeader

	if len(t.content) == 0 {
		return
	}

	rows := len(t.content[0])
	cols := len(t.content)

	trans := make([][]string, rows)

	for i := 0; i < rows; i++ {
		trans[i] = make([]string, cols)
		for j := 0; j < cols; j++ {
			trans[i][j] = t.content[j][i]
		}
	}

	t.content = trans
}

func (t *T) AddColHeader(c []string) error {
	if len(t.content) != 0 && len(t.content) != len(c) {

		fmt.Println(len(t.content))
		return rowSizeErr
	}
	h := make([]string, len(c))
	copy(h, c)
	t.colHeader = h
	return nil
}

func (t *T) AddWithHeader(vs ...[]string) error {
	h := make([]string, 0)
	c := make([][]string, len(vs))
	for i, v := range vs {
		if len(t.content) != 0 && len(t.content[0]) != len(v)-1 {
			return rowSizeErr
		}
		h = append(h, v[0])
		c[i] = v[1:]
	}

	t.content = append(t.content, c...)
	t.rowHeader = append(t.rowHeader, h...)
	if len(t.content) != 0 && len(t.colHeader) < len(t.content[0]) {
		t.colHeader = append(t.colHeader,
			make([]string, len(t.content[0])-len(t.colHeader))...)
	}
	return nil
}

func (t T) GetContent() [][]string {
	out := make([][]string, 0)
	for _, r := range t.content {
		newRow := make([]string, len(r))
		copy(newRow, r)
		out = append(out, newRow)
	}
	return out
}

func (t T) GetContentWithColHeader() [][]string {
	out := make([][]string, 0)
	out = append(out, t.colHeader)
	for _, c := range t.GetContent() {
		out = append(out, c)
	}
	return out
}

func prepend(base [][]string, add []string) [][]string {
	out := make([][]string, 0)
	for i, l := range base {
		out = append(out, append([]string{add[i]}, l...))
	}
	return out
}

func (t T) GetContentWithRowHeader() [][]string {
	return prepend(t.content, t.rowHeader)
}

func (t T) GetContentWithHeaders() [][]string {
	out := make([][]string, 0)
	out = append(out, append([]string{""}, t.colHeader...))
	inter := prepend(t.content, t.rowHeader)

	for _, it := range inter {
		out = append(out, it)
	}

	return out
}
