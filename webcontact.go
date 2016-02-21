package nlgids

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webcontact"
)

// WebBooking sends an email to the recipients from a contact form from the website.
func (n *NLgids) WebContact(w http.ResponseWriter, r *http.Request) (int, error) {
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	phone, persons := r.PostFormValue("phone"), r.PostFormValue("persons")
	message := r.PostFormValue("message")

	if name == "" || email == "" || message == "" {
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

	contact := &webcontact.Contact{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Persons: persons,
		Message: message,
	}
	return sendContactMail(contact, n.Config.Recipients)
}

func sendContactMail(c *webcontact.Contact, rcpts []string) (int, error) {
	subject := c.MailSubject()
	body, err := c.MailBody()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	mail := email.NewContact(subject, body)
	if err := mail.Do(rcpts); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
