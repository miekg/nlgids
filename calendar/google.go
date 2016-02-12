package calendar

import (
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
)

// FreeBusy will retrieve all evens for this Calendar and mark each day as either free
// or busy depending on the All-Day events in the Google Calendar.
func (c *Calendar) FreeBusy() error {
	b, err := ioutil.ReadFile("/home/miek/NLgids-fcbeb7928cdb.json")
	if err != nil {
		return err
	}

	config, err := google.JWTConfigFromJSON(b, gcal.CalendarReadonlyScope)
	if err != nil {
		return err
	}
	config.Subject = "miek@miek.nl" // TODO: ans
	client := config.Client(oauth2.NoContext)

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
		println(i.Summary)
		if i.Start.DateTime != "" {
			continue
		}
		whenTime, _ := time.Parse("2006-01-02", when)
		if _, ok := c.days[whenTime]; ok {
			println("setting", when)
			c.days[whenTime] = busy
		}
	}
	return nil
}