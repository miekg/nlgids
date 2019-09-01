package webtocht

import (
	"bytes"
	"text/template"
)

// TochtMail is a customer tocht form.
type TochtMail struct {
	Name    string
	Kenmerk string
}

const templ = `Hallo Ans,

Dit is het bevestigings formulier voor "{{.Name}}". Met kenmerk "{{.Kenmerk}}".

Met vriendelijke groet,
    NLgids mailer
`

func (t *Tocht) MailBody() (*bytes.Buffer, error) {
	te := template.New("Tocht template")
	te, err := te.Parse(templ)
	if err != nil {
		return nil, err
	}

	tm := &TochtMail{Name: t.FullName, Kenmerk: t.Kenmerk}
	buf := &bytes.Buffer{}
	if err := te.Execute(buf, tm); err != nil {
		return nil, err
	}
	return buf, nil
}

func (t *Tocht) MailSubject() string {
	subject := "Bevestiging (" + t.Kenmerk + "): \"" + t.FullName + "\""
	return subject
}
