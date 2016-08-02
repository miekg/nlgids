package nlgids

import (
	"os"
	"testing"
	"time"

	"github.com/miekg/nlgids/calendar"
)

func TestCalendarNow(t *testing.T) {
	c, err := calendar.New("", "", "")
	if err != nil {
		t.Fatal(err)
	}
	c.FreeBusy()
	t.Log(c.HTML())
}

func TestCalendarHistoric(t *testing.T) {
	c, err := calendar.New("2015-12-01", "", "")
	if err != nil {
		t.Fatal(err)
	}
	c.FreeBusy()
	t.Log(c.HTML())
}

func TestFreeBusy(t *testing.T) {
	date := "2016-08-02" // planned holiday date
	subject := "ans@nlgids.london"
	secret := "/etc/caddy/NLgids-fcbeb7928cdb.json"

	if _, err := os.Open(secret); err != nil {
		t.Logf("can open secret file, not performing test: %s", err)
		return
	}

	c, err := calendar.New(date, subject, secret)
	if err != nil {
		t.Errorf("can get new calendar: %s", err)
	}
	c.FreeBusy()

	tm := time.Date(2016, time.August, 2, 0, 0, 0, 0, time.UTC)
	days := c.Days()
	if days[tm].Available != calendar.Busy {
		t.Errorf("day %s should be busy, it is not", days[tm].Available)
	}
}
