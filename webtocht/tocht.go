package webtocht

import (
	"bytes"
	"io/ioutil"
	"strings"
	"text/template"

	"github.com/miekg/nlgids/date"
	"github.com/miekg/nlgids/ecb"
)

// Tocht holds all the data we need to generate a tocht.
type Tocht struct {
	FileName string // Name of the generated PDF.

	Kenmerk string // Unique kenmerk of this invoice.

	Tour     string
	Cost     float64
	Rate     float64 // current GBP:EUR rate, autofill
	Date     string  // YYYY/MM/DD form
	Name     string
	FullName string
	Email    string
}

// escapeTeX escapes TeX control characters. Currently: &, \ and %.
func escapeTeX(s string) string {
	s = strings.Replace(s, `&`, `\&`, -1)
	s = strings.Replace(s, `%`, `\%`, -1)
	s = strings.Replace(s, `_`, `\_`, -1)
	return s
}

func (t *Tocht) escapeTex() {
	t.Name = escapeTeX(t.Name)
	t.FullName = escapeTeX(t.FullName)
	t.Tour = escapeTeX(t.Tour)

	t.Name = strings.TrimSpace(t.Name)
	t.FullName = strings.TrimSpace(t.FullName)
	t.Tour = strings.TrimSpace(t.Tour)
}

// FillOut fills in these missing fields in i, such as:
// Day, Filename and makes Date Dutch.
func (t *Tocht) FillOut() (err error) {
	if t.FileName, err = date.NL("--date", t.Date, "+bevestiging-%d-%B-%Y.pdf"); err != nil {
		return err
	}
	if t.Date, err = date.NL("--date", t.Date, "+%d %B %Y"); err != nil {
		return err
	}
	if t.Rate, err = ecb.RateGBP(); err != nil {
		return err
	}
	return
}

// ExecuteTemplateAndWrite executes the template and writes the buffer to dst.
func (t *Tocht) ExecuteTemplateAndWrite(te *template.Template, name, dst string) error {
	buf := &bytes.Buffer{}
	if err := te.ExecuteTemplate(buf, name, t); err != nil {
		return err
	}

	if err := ioutil.WriteFile(dst, buf.Bytes(), 0664); err != nil {
		return err
	}
	return nil
}
