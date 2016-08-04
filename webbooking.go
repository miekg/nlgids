package nlgids

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/miekg/nlgids/email"
	ntour "github.com/miekg/nlgids/tour"
	"github.com/miekg/nlgids/webbooking"
)

// WebBooking sends an email to the recipients with a booking made from the website.
func (n *NLgids) WebBooking(w http.ResponseWriter, r *http.Request) (int, error) {
	tour, date := r.PostFormValue("tour"), r.PostFormValue("date")
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	phone, persons := r.PostFormValue("phone"), r.PostFormValue("persons")
	message := r.PostFormValue("message")

	if name == "" || email == "" {
		return http.StatusBadRequest, fmt.Errorf("nlgids: all empty")
	}
	if !strings.Contains(email, "@") {
		return http.StatusBadRequest, fmt.Errorf("nlgids: invalid email")
	}
	if persons != "" {
		if _, err := strconv.Atoi(persons); err != nil {
			return http.StatusBadRequest, err
		}
	}
	// Validate date and return error if not available.
	// use n.Config.Subject, n.Config.Secret

	tour = ntour.NameOrNonExists(tour, n.Config.Tours)

	booking := &webbooking.Booking{
		Tour:    tour,
		Date:    date,
		Name:    name,
		Email:   email,
		Phone:   phone,
		Persons: persons,
		Message: message,
	}
	log.Printf("[INFO] NLgids booking: %s:%s:%s:%s:%s:%s:%s", tour, date, name, email, phone, persons, message)

	return sendBookingMail(booking, n.Config.Recipients)
}

func sendBookingMail(b *webbooking.Booking, rcpts []string) (int, error) {
	subject := b.MailSubject()
	body, err := b.MailBody()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	mail := email.NewBooking(b.Email, subject, body)
	if err := mail.Do(rcpts); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
