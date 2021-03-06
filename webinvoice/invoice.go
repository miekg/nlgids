package webinvoice

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/miekg/nlgids/date"
	"github.com/miekg/nlgids/ecb"
)

// Invoice holds all the data we need to generate an invoice
type Invoice struct {
	FileName string // Name of the generated PDF.

	Kenmerk string // Unique kenmerk of this invoice.

	Tour     string
	Persons  int
	Time     string
	Duration string // 2:00
	Cost     float64
	Date     string // YYYY/MM/DD form
	Name     string
	FullName string
	Email    string // Has become optional.
	Where    string // Where to pickup.
	How      string // Ends in "om".

	Rate float64 // current GBP:EUR rate, autofill
	Day  string  // autofill

	OrigDate string // original values of Time and Data (before we converted to NL).
	OrigTime string
}

// escapeTeX escapes TeX control characters. Currently: &, \ and %.
func escapeTeX(s string) string {
	s = strings.Replace(s, `&`, `\&`, -1)
	s = strings.Replace(s, `%`, `\%`, -1)
	s = strings.Replace(s, `_`, `\_`, -1)
	return s
}

func (i *Invoice) escapeTex() {
	i.Name = escapeTeX(i.Name)
	i.FullName = escapeTeX(i.FullName)
	i.Tour = escapeTeX(i.Tour)
	i.Where = escapeTeX(i.Where)
	i.How = escapeTeX(i.How)

	i.Name = strings.TrimSpace(i.Name)
	i.FullName = strings.TrimSpace(i.FullName)
	i.Tour = strings.TrimSpace(i.Tour)
	i.Where = strings.TrimSpace(i.Where)
	i.How = strings.TrimSpace(i.How)
}

// FillOut fills in these missing fields in i, such as:
// Rate, Day, Filename and makes Date Dutch.
func (i *Invoice) FillOut() (err error) {
	i.OrigDate = i.Date
	i.OrigTime = i.Time
	if i.Rate, err = ecb.RateGBP(); err != nil {
		return err
	}
	if i.Day, err = date.NL("--date", i.Date, "+%A"); err != nil {
		return err
	}
	if i.FileName, err = date.NL("--date", i.Date, "+reservering-%d-%B-%Y.pdf"); err != nil {
		return err
	}
	if i.Date, err = date.NL("--date", i.Date, "+%d %B %Y"); err != nil {
		return err
	}
	return
}

// ExecuteTemplateAndWrite executes the template and writes the buffer to dst.
func (i *Invoice) ExecuteTemplateAndWrite(t *template.Template, name, dst string) error {
	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, name, i); err != nil {
		return err
	}

	if err := ioutil.WriteFile(dst, buf.Bytes(), 0664); err != nil {
		return err
	}
	return nil
}
