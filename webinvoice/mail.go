package webinvoice

import (
	"bytes"
	"text/template"
)

// Invoice is a customer invoice form.
type InvoiceMail struct {
	Name    string
	Kenmerk string
}

const templ = `Hallo Ans,

Dit is het reserverings formulier voor "{{.Name}}". Met kenmerk ""{{.Kenmerk}}".

Met vriendelijke groet,
    NLgids mailer
`

func (i *Invoice) MailBody() (*bytes.Buffer, error) {
	t := template.New("Invoice template")
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
	}

	c := &InvoiceMail{Name: i.FullName, Kenmerk: i.Kenmerk}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, c); err != nil {
		return nil, err
	}
	return buf, nil
}

func (i *Invoice) MailSubject() string {
	subject := "Formulier (" + i.Kenmerk + "): \"" + i.FullName + "\""
	return subject
}
