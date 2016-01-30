package nlgids

import (
	"testing"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webcontact"
)

func newContact() *webcontact.Contact {
	return &webcontact.Contact{
		Name:    "Miek Gieben",
		Email:   "miek@miek.nl",
		Phone:   "07774 517 566",
		Message: "Hee, hoe is het daar?",
	}
}

func TestContactCreate(t *testing.T) {
	c := newContact()
	subject := c.MailSubject()
	body, err := c.MailBody()
	if err != nil {
		t.Fatal(err)
	}
	mail := email.NewContact(subject, body)

	if mail.Subject != "[NLgids] Contact van \"Miek Gieben\"" {
		t.Fatal("wrong email Subject")
	}
	if mail.From != "nlgids@nlgids.london" {
		t.Fatal("wrong email From")
	}
	if len(mail.Cc) != 0 {
		t.Fatal("wrong email Cc")
	}
	if err := mail.Do(); err != nil {
		t.Fatalf("can't send mail %s: ", err)
	}
}
