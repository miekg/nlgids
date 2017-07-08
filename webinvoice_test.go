package nlgids

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/tour"
	"github.com/miekg/nlgids/webinvoice"
)

const TestFileName = "reservering-10-december-2015.pdf"

func newInvoice() *webinvoice.Invoice {
	i := &webinvoice.Invoice{
		Tour:     "walks/koninklijke",
		Persons:  2,
		Time:     "11.00",
		Duration: "2.0",
		Cost:     50.0,
		Date:     "2015/12/10",
		Name:     "Christel",
		FullName: "Christel Achternaam",
		Email:    "christel_bla@miek.nl",
		Where:    "Green Park metro station",
		How:      "Ik sta buiten de de fontein om",
		Kenmerk:  webinvoice.Kenmerk(time.Now().UTC()),
	}
	i.Tour = tour.NameOrNonExists(i.Tour, "/home/miek/html/nlgids.london/tours.json")
	return i
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
	if i.FileName != TestFileName {
		t.Fatalf("filename should be %s, got %s", TestFileName, i.FileName)
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
	ioutil.WriteFile("/tmp/test.pdf", pdf, 0644)

	body, err := i.MailBody()
	if err != nil {
		t.Fatal(err)
	}
	mail := email.NewInvoice(i.MailSubject(), body, i.FileName, pdf)
	if strings.Contains("Christel Achternaam", mail.Subject) {
		t.Fatal("wrong email Subject")
	}
	if mail.From != "" { // Set in mail.Do()
		t.Fatal("wrong email From")
	}
	if len(mail.Cc) != 0 {
		t.Fatal("wrong email Cc")
	}
	if len(mail.Attachments) != 1 {
		t.Fatal("wrong email attachment number")
	}
}
