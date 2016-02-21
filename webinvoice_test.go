package nlgids

import (
	"path"
	"testing"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webinvoice"
)

func newInvoice() *webinvoice.Invoice {
	return &webinvoice.Invoice{
		Tour:     "Van Koninklijke Huize",
		Persons:  2,
		Time:     "11.00",
		Duration: 2.0,
		Cost:     50.0,
		Date:     "2015/12/10",
		Name:     "Christel",
		FullName: "Christel Achternaam",
		Email:    "christel@miek.nl",
		Where:    "Green Park metro station",
		How:      "Ik sta buiten de de fontein om",
	}
}

func TestInvoiceFill(t *testing.T) {
	i := newInvoice()
	err := i.FillOut()
	if err != nil {
		t.Fatal(err)
	}
	if i.Day == "" {
		t.Fatal("should be non empty: 'Day'")
	}
	if i.FileName == "" {
		t.Fatal("should be non empty: 'FileName'")
	}
	if i.Rate == 0 {
		t.Fatal("should be non empty: 'Rate'")
	}
}

func TestTexFiles(t *testing.T) {
	sources, err := webinvoice.TeXFiles(templateDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(sources) == 0 {
		t.Fatalf("no sources found in %s", templateDir)
	}
}

func TestInvoiceCreate(t *testing.T) {
	invoice := path.Join(templateDir, webinvoice.Template)

	i := newInvoice()
	pdf, err := i.Create(templateDir, invoice)
	if err != nil {
		t.Fatal(err)
	}
	if len(pdf) == 0 {
		t.Fatal("no pdf produced")
	}
	body, err := i.MailBody()
	if err != nil {
		t.Fatal(err)
	}
	mail := email.NewInvoice(i.MailSubject(), body, i.FileName, pdf)
	if mail.Subject != "[NLgids] Formulier \"Christel Achternaam\"" {
		t.Fatal("wrong email Subject")
	}
	if mail.From != "" {  // Set in mail.Do()
		t.Fatal("wrong email From")
	}
	if len(mail.Cc) != 0 {
		t.Fatal("wrong email Cc")
	}
	if len(mail.Attachments) != 1 {
		t.Fatal("wrong email attachment number")
	}
}
