package nlgids

import (
	"fmt"
	"net/http"

	"github.com/miekg/nlgids/calendar"
)

// WebCalendar returns a calendar in table form. All-day events from the
// subject are greyed out, as are pasted days.
func (n *NLgids) WebCalendar(w http.ResponseWriter, r *http.Request) (int, error) {
	date := r.PostFormValue("date")     // YYYY-MM-DD, empty is allowed
	tour := r.PostFormValue("tourtype") // see tours.json in the site, this is the "type"
	tour := r.PostFormValue("tourid")   // see tours.json in the site, this is the "id"

	c, err := calendar.New(date, n.Config.Subject, n.Config.Secret)
	if err != nil {
		return http.StatusBadRequest, nil
	}
	c.FreeBusy()
	fmt.Fprintf(w, c.HTML())

	return http.StatusOK, nil
}
