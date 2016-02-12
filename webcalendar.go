package nlgids

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/miekg/nlgids/calendar"
)

func TestCalendar(t *testing.T) {
	now := time.Now().UTC()
	c, err := calendar.New(now.Format("2006-01-02"))
	if err != nil {
		t.Fatal(err)
	}
	c.FreeBusy()
	t.Log(c.HTML())
}

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
