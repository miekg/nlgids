package nlgids

import (
	"fmt"
	"net/http"
	"time"

	"github.com/miekg/nlgids/calendar"
)

func WebCalendarTest(w http.ResponseWriter, r *http.Request) (int, error) {
	date := r.PostFormValue("date") // YYYY-MM-DD
	if date == "" {
		return http.StatusBadRequest, nil
	}
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	c, err := calendar.New(date)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	c.FreeBusy()
	fmt.Fprintf(w, c.HTML())

	return http.StatusOK, nil
}
