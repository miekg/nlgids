package webbooking

// not ready

import (
	"bytes"
	"text/template"
)

// Booking is a customer booking form.
type Booking struct {
	Tour string
	Date string

	Name    string
	Email   string
	Phone   string
	Persons string
	Message string
}

const templ = `Hallo Ans,

Er is een boekings formulier ingevuld, met de volgende details:

* Tour......: {{.Tour}} op {{.Date}}
* Datum.....: {{.Date}}
* Naam......: {{.Name}} ({{.Email}})
* Tel.......: {{.Phone}}
* Personen..: {{.Persons}}

En het volgende bericht is achter gelaten:

=====================================================================

{{.Message}}

=====================================================================

Met vriendelijke groet,
    NLgids mailer
`

func (b *Booking) MailBody() (*bytes.Buffer, error) {
	t := template.New("Contact template")
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, b); err != nil {
		return nil, err
	}
	return buf, nil
}

func (b *Booking) MailSubject() string {
	subject := ""
	if b.Date == "" {
		subject = "Boeking van \"" + b.Name + "\""
	} else {
		subject = "Boeking op " + b.Date + " van \"" + b.Name + "\""
	}
	return subject
}
