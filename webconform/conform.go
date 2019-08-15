package webconform

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/miekg/nlgids/date"
	"github.com/miekg/nlgids/ecb"
)

// Conform holds all the data we need to generate a conform
type Conform struct {
	FileName string // Name of the generated PDF.

	Tour     string
	Persons  int
	Time     string
	Duration string // 2:00
	Cost     float64
	Rate     float64 // current GBP:EUR rate, autofill
	Date     string  // YYYY/MM/DD form
	Name     string
	FullName string
	Email    string // Has become optional.
	Where    string // Where to pickup.
	How      string // Ends in "om".

	Day string // autofill
}

// escapeTeX escapes TeX control characters. Currently: &, \ and %.
func escapeTeX(s string) string {
	s = strings.Replace(s, `&`, `\&`, -1)
	s = strings.Replace(s, `%`, `\%`, -1)
	s = strings.Replace(s, `_`, `\_`, -1)
	return s
}

func (c *Conform) escapeTex() {
	c.Name = escapeTeX(c.Name)
	c.FullName = escapeTeX(c.FullName)
	c.Tour = escapeTeX(c.Tour)
	c.Where = escapeTeX(c.Where)
	c.How = escapeTeX(c.How)

	c.Name = strings.TrimSpace(c.Name)
	c.FullName = strings.TrimSpace(c.FullName)
	c.Tour = strings.TrimSpace(c.Tour)
	c.Where = strings.TrimSpace(c.Where)
	c.How = strings.TrimSpace(c.How)
}

// FillOut fills in these missing fields in i, such as:
// Day, Filename and makes Date Dutch.
func (c *Conform) FillOut() (err error) {
	if c.Day, err = date.NL("--date", c.Date, "+%A"); err != nil {
		return err
	}
	if c.FileName, err = date.NL("--date", c.Date, "+bevestiging-%d-%B-%Y.pdf"); err != nil {
		return err
	}
	if c.Date, err = date.NL("--date", c.Date, "+%d %B %Y"); err != nil {
		return err
	}
	if c.Rate, err = ecb.RateGBP(); err != nil {
		return err
	}
	return
}

// ExecuteTemplateAndWrite executes the template and writes the buffer to dst.
func (c *Conform) ExecuteTemplateAndWrite(t *template.Template, name, dst string) error {
	buf := &bytes.Buffer{}
	if err := t.ExecuteTemplate(buf, name, c); err != nil {
		return err
	}

	if err := ioutil.WriteFile(dst, buf.Bytes(), 0664); err != nil {
		return err
	}
	return nil
}
