package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

var (
	csvSuff    = ".csv"
	timeFormat = "2006-02-01"
)

func GetConfigFromFile(filename string, config interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Cannot read file %q", filename)
	}
	return GetConfig(b, config)
}

func GetConfig(b []byte, config interface{}) error {
	return yaml.Unmarshal(b, config)
}

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

type Filename struct {
	Date                time.Time
	Stock, Act, Id, Ext string
}

func (f Filename) String() string {
	s := f.Date.Format(timeFormat) + "-" + f.Stock + "-" + f.Id + "-" + f.Act + f.Ext
	return s
}
