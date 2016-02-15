package nlgids

import (
	"fmt"
	"net/http"

	"github.com/miekg/nlgids/calendar"
)

func WebCalendarTest(w http.ResponseWriter, r *http.Request) (int, error) {
	date := r.PostFormValue("date") // YYYY-MM-DD, empty is allowed.
	c, err := calendar.New(date)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	c.FreeBusy()
	fmt.Fprintf(w, c.HTML())

	return http.StatusOK, nil
}
