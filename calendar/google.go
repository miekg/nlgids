package calendar

import (
	"fmt"
	"io/ioutil"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcal "google.golang.org/api/calendar/v3"
)

// FreeBusy will retrieve all evens for this Calendar and mark each day as either free
// or busy depending on the All-Day events in the Google Calendar.
func (c *Calendar) FreeBusy() error {
	b, err := ioutil.ReadFile("/home/miek/downloads/NLgids-fcbeb7928cdb.json")
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

	fmt.Println("Upcoming events:")
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			// for each when and all-day events we set the calendar entry to busy
			// (if not in the past).
			fmt.Printf("%s (%s)\n", i.Summary, when)
		}
	} else {
		fmt.Printf("No upcoming events found.\n")
	}
	return nil
}
