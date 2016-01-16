package main

import (
	"fmt"
	. "github.com/olivier5741/stock-manager/item"
	"path/filepath"
	"strings"
	"time"
)

var (
	csvSuff    = ".csv"
	timeFormat = "2006-02-01"
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
	if l := len(args); l != 6 {
		err = fmt.Errorf("There should be 6 hyphen seperated args in %q", s)
		return
	}

	//Date
	date, err := time.Parse(timeFormat, args[0]+"-"+args[1]+"-"+args[2])
	if err != nil {
		return
	}
	f.Date = date

	//Stock
	f.Stock = args[3]

	//Id
	f.Id = args[4]

	//Action
	f.Act = args[5]

	return
}

type Config struct {
	Stock  []string
	Redist []ConfigRedist
	Prod   []ConfigProd
}

type ConfigRedist struct {
	From string
	To   []string
}

type ConfigDispo struct {
	Ok  bool
	Pmg []string
}

type ConfigProd struct {
	Id, Name, Unit  string
	Bulk, Room, Min int
	//	Dispo           ConfigDispo
}

func (c Config) GetMissingItems() (out Items) {
	out = map[string]Item{}
	for _, p := range c.Prod {
		out[p.Name] = Item{Prod(p.Name), Val{p.Min * p.Bulk}}
	}
	return
}

type Filename struct {
	Date                time.Time
	Stock, Act, Id, Ext string
}

func (f Filename) String() string {
	s := f.Date.Format(timeFormat) + "-" + f.Stock + "-" + f.Id + "-" + f.Act + f.Ext
	return s
}
