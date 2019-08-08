package nlgids

import (
	"fmt"
	"net/http"

	"github.com/miekg/nlgids/calendar"
)

// WebCalendar returns a calendar in table form. All-day events from the
// subject are greyed out, as are pasted days.
func (n *NLgids) WebCalendar(w http.ResponseWriter, r *http.Request) (int, error) {
	date := r.PostFormValue("date") // YYYY-MM-DD, empty is allowed

	c, err := calendar.New(date, n.Config.Subject, n.Config.Secret)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	c.FreeBusy()
	fmt.Fprintf(w, c.HTML())

	return http.StatusOK, nil
}
