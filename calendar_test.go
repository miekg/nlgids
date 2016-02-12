package nlgids

import (
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
	println(c.HTML())
}
