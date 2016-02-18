package sheet

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	sep = "-"
	timeFormat = "2006" + sep + "01" + sep + "02"
	numberPrefix    = "n"
	Tr = func(s string) string {return s}
)

// sample en attente-2016-02-18-n1-entree
type Filename struct {
	Status  string
	Date    time.Time
	ID, Act string
	Basic string
}

func NewBasicFilename(s string) Filename {
	return Filename{Basic: generatedPrefix + s}
}

func NewFilename(s string) (Filename, error) {
	var status string
	args := strings.Split(s, sep)

	if l := len(args); l != 5 && l != 6 {
		return Filename{}, fmt.Errorf(Tr("filename_number_argument_error"), s)
	}

	if len(args) == 6 {
		status = args[0]
		args = args[1:]
	}

	date, err := time.Parse(timeFormat, args[0]+sep+args[1]+sep+args[2])
	if err != nil {
		return Filename{}, fmt.Errorf(Tr("filename_date_parse_error"), s)
	}

	return Filename{status, date, strings.TrimPrefix(args[3], numberPrefix), args[4], ""}, nil
}

func NewDraftFilename(id int, act string) Filename {
	return Filename{Tr("file_prefix_draft"), time.Now(), strconv.Itoa(id), act, ""}
}

func (f Filename) String() string {
	if f.Basic != "" {
		return f.Basic
	} 

	var s string
	if f.Status != "" {
		s = f.Status + sep + s
	}
	return s + f.Date.Format(timeFormat) + sep + numberPrefix + f.ID + sep + f.Act
}

func (f Filename) Time() string {
	return f.Date.Format(timeFormat) + sep + numberPrefix + f.ID
}