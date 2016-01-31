package nlgids

import (
	"log"
	"net/http"

	"github.com/miekg/nlgids/email"
	"github.com/miekg/nlgids/webcontact"
)

func WebContactTest(w http.ResponseWriter, r *http.Request) (int, error) {
	testContact := &webcontact.Contact{
		Name:    "Miek Gieben",
		Email:   "miek@miek.nl",
		Phone:   "07774 517 566",
		Message: "Hee, hoe is het daar?",
	}
	return sendContactMail(testContact)
}

func WebContact(w http.ResponseWriter, r *http.Request) (int, error) {
	name, email := r.PostFormValue("name"), r.PostFormValue("email")
	phone, message := r.PostFormValue("phone"), r.PostFormValue("message")
	// validate and return http.StatusBadRequest

	contact := &webcontact.Contact{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Message: message,
	}
	return sendContactMail(contact)
}

func sendContactMail(c *webcontact.Contact) (int, error) {
	subject := c.MailSubject()
	body, err := c.MailBody()
	if err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	mail := email.NewContact(subject, body)
	if err := mail.Do(); err != nil {
		log.Printf("%s", err)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
