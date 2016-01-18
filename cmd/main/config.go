package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	. "github.com/olivier5741/stock-manager/item"
	"path/filepath"
	"strings"
	"time"
)

var (
	csvSuff    = ".csv"
	timeFormat = "2006-01-02"
)

func ParseFilename(s string) (f Filename, err error) {
	f = Filename{}

	s = filepath.Base(s)
	//Extension
	f.Ext = filepath.Ext(s)
	if strings.HasSuffix(s, csvSuff) {
		err = fmt.Errorf("No csv suffix in filename: %q", s)
	}

	ts := strings.TrimSuffix(s, f.Ext)

	args := strings.Split(ts, "-")
	if l := len(args); l != 5 && l != 6 {
		err = fmt.Errorf("There should be 5 or 6 hyphen seperated args in %q", s)
		return
	}

	//Date
	date, err := time.Parse(timeFormat, args[0]+"-"+args[1]+"-"+args[2])
	if err != nil {
		return
	}
	f.Date = date

	//Stock
	//f.Stock = args[3]
	f.Stock = "bievre" //To refactor

	//Id
	f.Id = strings.TrimPrefix(args[3], "n°")

	//Action
	f.Act = args[4]

	if len(args) == 6 {
		f.Status = args[5]
	}

	return
}

type ConfigProd struct {
	Id, Name  string
	Bulk, Min int
}

func GetMissingItems(c []ConfigProd) (out Items) {
	out = map[string]Item{}
	for _, p := range c {
		out[p.Name] = Item{Prod(p.Name), Val{p.Min * p.Bulk}}

		log.Debug(p)
	}
	return
}

func ToItemStringLines(prods []ConfigProd) (out [][]string) {
	out = make([][]string, len(prods))
	for i, prod := range prods {
		out[i] = []string{prod.Name, ""}
	}
	return
}

type Filename struct {
	Date                        time.Time
	Stock, Act, Id, Ext, Status string
}

func (f Filename) String() string {
	s := f.Date.Format(timeFormat) + "-n°" + f.Id + "-" + f.Act
	if f.Status != "" {
		s = s + "-" + f.Status
	}
	s = s + f.Ext

	return s
}
