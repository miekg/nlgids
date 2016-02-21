package nlgids

import (
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webinvoice"
)

// WebInvoice sends an email to the recipients from a invoice form from the website.
func (n *NLgids) WebInvoice(w http.ResponseWriter, r *http.Request) (int, error) {
	tour, personsStr := r.PostFormValue("tour"), r.PostFormValue("persons")
	time, durationStr := r.PostFormValue("time"), r.PostFormValue("duration")
	costStr, date := r.PostFormValue("cost"), r.PostFormValue("date")
	name, fullname := r.PostFormValue("name"), r.PostFormValue("fullname")
	email := r.PostFormValue("email")
	where, how := r.PostFormValue("where"), r.PostFormValue("how")

	if tour == "" || personsStr == "" || time == "" || date == "" || name == "" ||
		fullname == "" || email == "" {
		return http.StatusBadRequest, nil
	}

	if !strings.Contains(email, "@") {
		return http.StatusBadRequest, nil
	}
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	persons, err := strconv.Atoi(personsStr)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	invoice := &webinvoice.Invoice{
		Tour:     tour,
		Persons:  persons,
		Time:     time,
		Duration: duration,
		Cost:     cost,
		Date:     date,
		Name:     name,
		FullName: fullname,
		Email:    email,
		Where:    where,
		How:      how,
	}

	tmpl := path.Join(n.Config.Template, webinvoice.Template)

	pdf, err := invoice.Create(n.Config.Template, tmpl)
	if err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	if len(pdf) == 0 {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return sendInvoice(invoice, pdf, n.Config.Recipients)
}

func sendInvoice(i *webinvoice.Invoice, pdf []byte, rcpts []string) (int, error) {
	subject := i.MailSubject()
	body, err := i.MailBody()
	if err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}

	mail := email.NewInvoice(subject, body, i.FileName, pdf)
	if err := mail.Do(rcpts); err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
