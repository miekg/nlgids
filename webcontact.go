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
	subject := testContact.MailSubject()
	body, err := testContact.MailBody()
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
