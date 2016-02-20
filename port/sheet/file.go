package sheet

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/olivier5741/stock-manager/asset"
)

var (
	sep = "-"
	timeFormat = "2006" + sep + "01" + sep + "02"
	numberPrefix    = "n"
	Tr = func(s string) string {return s}
)

type BasicFilename struct {
	Name, UID string
}

// sample en attente-2016-02-18-n1-entree
type Filename struct {
	Status  string
	Date    time.Time
	ID, Act string
	Bypass string
	UID string
}

func NewBasicFilename(s string) Filename {
	return Filename{Bypass: generatedPrefix + s}
}

func NewFilename(s BasicFilename) (Filename, error) {
	var status string
	args := strings.Split(s.Name, sep)

	if l := len(args); l != 5 && l != 6 {
		return Filename{}, fmt.Errorf(asset.Tr("filename_number_argument_error"), s)
	}

	if len(args) == 6 {
		status = args[0]
		args = args[1:]
	}

	date, err := time.Parse(timeFormat, args[0]+sep+args[1]+sep+args[2])
	if err != nil {
		return Filename{}, fmt.Errorf(asset.Tr("filename_date_parse_error"), s)
	}

	return Filename{status, date, strings.TrimPrefix(args[3], numberPrefix), args[4], "", s.UID}, nil
}

func NewDraftFilename(id int, act string) Filename {
	return Filename{asset.Tr("file_prefix_draft"), time.Now(), strconv.Itoa(id), act, "", ""}
}

func (f Filename) String() string {
	if f.Bypass != "" {
		return f.Bypass
	} 

	var s string
	if f.Status != "" {
		s = f.Status + sep + s
	}
	return s + f.Date.Format(timeFormat) + sep + numberPrefix + f.ID + sep + f.Act
}

func (f Filename) Basic() BasicFilename {
	return BasicFilename{f.String(),f.UID}
}

func (f Filename) Time() string {
	return f.Date.Format(timeFormat) + sep + numberPrefix + f.ID
}