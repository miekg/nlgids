package webcontact

import (
	"bytes"
	"text/template"
)

// ContactMail is a customer contact form.
type ContactMail struct {
	Name    string
	Email   string
	Phone   string
	Message string
}

const templ = `Hallo Ans,

Er is een contact formulier ingevuld, met de volgende details:

* Naam..: {{.Name}} ({{.Email}})
* Tel...: {{.Phone}}

En het volgende bericht is achter gelaten:

=====================================================================

{{.Message}}

=====================================================================

Met vriendelijke groet,
    NLgids mailer
`

func MailBody(name string) (*bytes.Buffer, error) {
	t := template.New("Contact template")
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
	}

	c := &ContactMail{
		"Miek Gieben",
		"miek@miek.nl",
		"07774 517 566",
		"Hee, we komen op bezoek!",
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, c); err != nil {
		return nil, err
	}
	return buf, nil
}

func MailSubject(name string) string {
	subject := "[NLgids] Contact van \""+name+"\""
	return subject
}
