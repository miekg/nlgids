package calendar

import (
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
)

// subject = "miek@miek.nl" // ans@nlgids.london
// secret  = "/etc/caddy/NLgids-fcbeb7928cdb.json"

func client(subject, secret string) (*http.Client, error) {
	b, err := ioutil.ReadFile(secret)
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(b, gcal.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	config.Subject = subject
	client := config.Client(oauth2.NoContext)
	return client, nil
}

// Diff. name
// FreeBusy returns true if there is an all-day event on the the date d (YYYY-MM-DD).
// func (c *Calendar) FreeBusy(d string) (bool, error) {
// // Check this one date
// return true, nil
// }

// FreeBusy will retrieve all evens for this Calendar and mark each day as either free
// or busy depending on the All-Day events in the Google Calendar.
// TODO(miek): This works only for 1-day freebusy events, not for multiday events!
func (c *Calendar) FreeBusy() error {
	client, err := client(c.subject, c.secret)
	if err != nil {
		return err
	}

	srv, err := gcal.New(client)
	if err != nil {
		return err
	}

	// TimeMax is exclusive, so we need to add another day to c.end to get all the events we want.
	begin := c.begin.Format(time.RFC3339)
	end := c.end.AddDate(0, 0, 1).Format(time.RFC3339)

	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(begin).TimeMax(end).OrderBy("startTime").Do()
	if err != nil {
		return err
	}

	for _, i := range events.Items {
		when := i.Start.Date
		// If the DateTime is an empty string the Event is an all-day Event.
		// So only Date is available.
		if i.Start.DateTime != "" {
			continue
		}
		whenTime, _ := time.Parse("2006-01-02", when)
		if _, ok := c.days[whenTime]; ok {
			c.days[whenTime] = AvailMeta{Available: Busy}
		}
	}
	return nil
}
