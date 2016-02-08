package nlgids

import (
	"log"
	"net/http"
	"path"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webinvoice"
)

func WebInvoiceTest(w http.ResponseWriter, r *http.Request) (int, error) {
	testInvoice := &webinvoice.Invoice{
		Tour:     "Van Koninklijke Huize",
		Persons:  2,
		Time:     "11.00",
		Duration: 2.0,
		Cost:     50.0,
		Date:     "2015/12/10",
		Name:     "Miek",
		FullName: "Miek Gieben",
		Email:    "miek@miek.nl",
		Where:    "Green Park metro station",
		How:      "Ik sta buiten de de fontein om",
	}

	tmpl := path.Join(templateDir, webinvoice.DefaultTemplate)

	pdf, err := testInvoice.Create(templateDir, tmpl)
	if err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	if len(pdf) == 0 {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return sendInvoice(testInvoice, pdf)
}

func sendInvoice(i *webinvoice.Invoice, pdf []byte) (int, error) {
	subject := i.MailSubject()
	body, err := i.MailBody()
	if err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}

	mail := email.NewInvoice(subject, body, i.FileName, pdf)
	if err := mail.Do(); err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
