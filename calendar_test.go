package nlgids

import (
	"testing"

	"github.com/miekg/nlgids/calendar"
)

func TestCalendarNow(t *testing.T) {
	c, err := calendar.New("")
	if err != nil {
		t.Fatal(err)
	}
	c.FreeBusy()
	t.Log(c.HTML())
}

func TestCalendarHistoric(t *testing.T) {
	c, err := calendar.New("2015-12-01")
	if err != nil {
		t.Fatal(err)
	}
	c.FreeBusy()
	t.Log(c.HTML())
}
