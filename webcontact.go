package nlgids

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webcontact"
)

func (n *NLgids) WebContact(w http.ResponseWriter, r *http.Request) (int, error) {
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	phone, persons := r.PostFormValue("phone"), r.PostFormValue("persons")
	message := r.PostFormValue("message")

	if name == "" || email == "" || message == "" {
		return http.StatusBadRequest, nil
	}
	if !strings.Contains(email, "@") {
		return http.StatusBadRequest, nil
	}
	if persons != "" {
		if _, err := strconv.Atoi(persons); err != nil {
			return http.StatusBadRequest, nil
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
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	mail := email.NewContact(subject, body)
	if err := mail.Do(rcpts); err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
