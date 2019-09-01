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
	"github.com/miekg/nlgids/webtocht"
)

// WebTocht sends an email to the recipients from a invoice form from the website.
func (n *NLgids) WebTocht(w http.ResponseWriter, r *http.Request) (int, error) {
	tour := r.PostFormValue("tour")
	costStr := r.PostFormValue("cost")
	name, fullname := r.PostFormValue("name"), r.PostFormValue("fullname")
	email := r.PostFormValue("email")
	date := r.PostFormValue("date")

	if date == "" || tour == "" || costStr == "" || name == "" || fullname == "" {
		return http.StatusBadRequest, fmt.Errorf("nlgids: all empty")
	}

	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Get the real the name of the tour.
	tour = ntour.NameOrNonExists(tour, n.Config.Tours)

	log.Printf("[INFO] NLgids tour: %s:%s:%s:%s:%f", tour, name, email, date, cost)

	tocht := &webtocht.Tocht{
		Tour:     tour,
		Cost:     cost,
		Name:     name,
		Date:     date,
		FullName: fullname,
		Email:    email,
		Kenmerk:  webtocht.Kenmerk(time.Now().UTC()),
	}

	tmpl := path.Join(n.Config.Template, webtocht.Template)

	pdf, err := tocht.Create(n.Config.Template, tmpl)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(pdf) == 0 {
		return http.StatusInternalServerError, err
	}
	return sendTocht(tocht, pdf, n.Config.Recipients)
}

func sendTocht(t *webtocht.Tocht, pdf []byte, rcpts []string) (int, error) {
	subject := t.MailSubject()
	body, err := t.MailBody()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	mail := email.NewTocht(subject, body, t.FileName, pdf)
	if err := mail.Do(rcpts); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
