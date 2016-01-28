package webbooking

// not ready

import (
	"bytes"
	"log"
	"text/template"

	"github.com/miekg/nlgids/email"
)

// Booking is a customer booking form.
type Contact struct {
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

func doit() {
	t := template.New("Contact template")
	t, err := t.Parse(templ)
	if err != nil {
		log.Fatal(err)
	}

	c := &Contact{
		"Miek Gieben",
		"miek@miek.nl",
		"07774 517 566",
		"Hee, we komen op bezoek!",
	}

	buf := &bytes.Buffer{}

	if err := t.Execute(buf, c); err != nil {
		log.Fatal(err)
	}
	mail := email.NewContactEmail("[NLgids] Contact van \""+c.Name+"\"", buf)
	if err := email.Do(mail); err != nil {
		log.Fatal(err)
	}
}
