package drivefile

import(
	"testing"
	"encoding/csv"
)

func TestGetAll(t *testing.T) {

	d := DriveFile{
		"0BzIZ3dfuz-CEN2dfQ1liU0x6eVU",
		GetService()}
	t.Log(d.GetAll())
}

func TestNewReader(t *testing.T) {
	h := Header{"1uW9g6IyUmJxAQ0hiZsTwQVwmUu4VcPdA5-g-TN4HI0E",	"no_name"}
	d := DriveFile{
		"0BzIZ3dfuz-CEN2dfQ1liU0x6eVU",
		GetService()}

	r := d.NewReader(h)
	defer r.Close()

	rcsv := csv.NewReader(r)
	out, _ := rcsv.ReadAll()
	t.Log(out)
}

func TestNewWriter(t *testing.T) {
	h := Header{"",	"test"}
	d := DriveFile{
		"0BzIZ3dfuz-CEN2dfQ1liU0x6eVU",
		GetService()}

	w := d.NewWriter(h)
	defer w.Close()

	wcsv := csv.NewWriter(w)
 	wcsv.WriteAll([][]string{[]string{"test"}})
}