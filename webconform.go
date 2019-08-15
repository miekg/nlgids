package nlgids

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/miekg/nlgids/email"
	ntour "github.com/miekg/nlgids/tour"
	"github.com/miekg/nlgids/webconform"
)

// WebConform sends an email to the recipients from a invoice form from the website.
func (n *NLgids) WebConform(w http.ResponseWriter, r *http.Request) (int, error) {
	tour, personsStr := r.PostFormValue("tour"), r.PostFormValue("persons")
	t, duration := r.PostFormValue("time"), r.PostFormValue("duration")
	costStr, date := r.PostFormValue("cost"), r.PostFormValue("date")
	name, fullname := r.PostFormValue("name"), r.PostFormValue("fullname")
	email := r.PostFormValue("email")
	where, how := r.PostFormValue("where"), r.PostFormValue("how")

	if tour == "" || personsStr == "" || t == "" || date == "" || name == "" ||
		fullname == "" || duration == "" {
		return http.StatusBadRequest, fmt.Errorf("nlgids: all empty")
	}

	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		return http.StatusBadRequest, err
	}
	persons, err := strconv.Atoi(personsStr)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Get the real the name of the tour.
	tour = ntour.NameOrNonExists(tour, n.Config.Tours)

	log.Printf("[INFO] NLgids tour: %s:%s:%s:%s:%s:%f", tour, name, email, personsStr, date, cost)

	conform := &webconform.Conform{
		Tour:     tour,
		Persons:  persons,
		Time:     t,
		Duration: duration,
		Cost:     cost,
		Date:     date,
		Name:     name,
		FullName: fullname,
		Email:    email,
		Where:    where,
		How:      how,
		Kenmerk:  webconform.Kenmerk(time.Now().UTC()),
	}

	tmpl := path.Join(n.Config.Template, webconform.Template)

	pdf, err := conform.Create(n.Config.Template, tmpl)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(pdf) == 0 {
		return http.StatusInternalServerError, err
	}
	return sendConform(conform, pdf, n.Config.Recipients)
}

func sendConform(c *webconform.Conform, pdf []byte, rcpts []string) (int, error) {
	subject := c.MailSubject()
	body, err := c.MailBody()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	mail := email.NewConform(subject, body, c.FileName, pdf)
	if err := mail.Do(rcpts); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
