package webconform

import (
	"bytes"
	"text/template"
)

// ConformMail is a customer conform form.
type ConformMail struct {
	Name string
}

const templ = `Hallo Ans,

Dit is het bevestigings formulier voor "{{.Name}}".

Met vriendelijke groet,
    NLgids mailer
`

func (c *Conform) MailBody() (*bytes.Buffer, error) {
	t := template.New("Conform template")
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
	}

	cm := &ConformMail{Name: c.FullName}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, cm); err != nil {
		return nil, err
	}
	return buf, nil
}

func (c *Conform) MailSubject() string {
	subject := "Bevestiging: \"" + c.FullName + "\""
	return subject
}
