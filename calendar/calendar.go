// Package calendar generates and has utlilty functions to generate an HTML calendar for use in nlgids.
package calendar

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"time"
)

var (
	avail  = [...]string{"", "past", "busy", "free"}
	meta   = [...]string{"", " now", " prev", " next"} // spaces before each word
	months = [...]string{"boogie", "januari", "februari", "maart", "april", "mei", "juni", "juli", "augustus", "september", "oktober", "november", "december"}
)

func (a Available) String() string { return avail[a] }
func (m Meta) String() string      { return meta[m] }

type (
	Available int
	Meta      int

	// The availabilty for a specific date.
	AvailMeta struct {
		Available
		Meta
	}
)

const (
	_ Available = iota
	Past
	Busy
	Free
)

const (
	_ Meta = iota
	Now
	Prev
	Next
)

const templ = `
    <div class="panel-heading text-center">
    <div class="row">
        <div class="col-md-1"> </div>

	<div class="col-md-10">
            <a class="btn btn-default btn-sm" onclick='BookingCalendar({{.Prev}})'>
                <span class="glyphicon glyphicon-arrow-left"></span>
            </a>

		&nbsp;<strong>{{.MonthNL}}</strong>&nbsp;

            <a class="btn btn-default btn-sm" onclick='BookingCalendar({{.Next}})'>
                <span class="glyphicon glyphicon-arrow-right"></span>
            </a>

	</div>

	<div class="col-md-1"> </div>
    </div>
</div>
`

type header struct {
	Prev    string
	Next    string
	MonthEN string
	MonthNL string
}

func Date(t time.Time) string {
	date := fmt.Sprintf("%4d-%02d-%02d", t.Year(), t.Month(), t.Day())
	return date
}

// Calendar holds the HTML that makes up the calendar. Each
// day is indexed by the 12 o' clock night time as a time.Time.
// All date are in the UTC timezone.
type Calendar struct {
	days  map[time.Time]AvailMeta
	begin time.Time
	end   time.Time
	start time.Time // generated for this date

	subject string // who's calendar
	secret  string // service account client_secret.json
}

type times []time.Time

func (t times) Len() int           { return len(t) }
func (t times) Less(i, j int) bool { return t[i].Before(t[j]) }
func (t times) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// Days returns the days of this calendar.
func (c *Calendar) Days() map[time.Time]AvailMeta { return c.days }

func (c *Calendar) heading() string {
	s := `<div class="row">
<div class="col-md-10 col-md-offset-1">
<table class="table table-bordered table-condensed">`
	s += "<tr><th>zo</th><th>ma</th><th>di</th><th>wo</th><th>do</th><th>vr</th><th>za</th></tr>\n"
	return s
}

// Header returns the header of the calendar.
func (c *Calendar) Header() (string, error) {
	month := c.start.Month()

	prev := c.start.AddDate(0, -1, 0)
	next := c.start.AddDate(0, +1, 0)
	monthEN := fmt.Sprintf("%s %d", month.String(), c.start.Year())
	monthNL := fmt.Sprintf("%s %d", months[month], c.start.Year())
	head := &header{
		Prev:    Date(prev),
		Next:    Date(next),
		MonthEN: monthEN,
		MonthNL: monthNL,
	}

	t := template.New("Header template")
	t, err := t.Parse(templ)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, head)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Calendar) Footer() string {
	return `</table>
</div>
</div>`
}

func (c *Calendar) openTR() string  { return "<tr>\n" }
func (c *Calendar) closeTR() string { return "</tr>\n" }

func (c *Calendar) entry(t time.Time) string {
	am := c.days[t]
	class := fmt.Sprintf("\t<td class=\"%s%s\">", am.Available, am.Meta)
	close := "</td>\n"
	href := ""

	switch am.Available {
	case Free:
		href = fmt.Sprintf("<a href=\"#\" onclick=\"BookingDate('%s')\">%d</a>", Date(t), t.Day()) // BookingDate is defined on the page/form itself
	case Busy, Past:
		href = fmt.Sprintf("%d", t.Day())
	}
	s := class + href + close
	return s
}

// HTML returns the calendar in a string containing HTML.
func (c *Calendar) HTML() string {
	s, _ := c.Header()
	s += c.html()
	s += c.Footer()
	return s
}

func (c *Calendar) sort() times {
	keys := times{}
	for k := range c.days {
		keys = append(keys, k)
	}

	sort.Sort(keys)
	return keys
}

func (c *Calendar) html() string {
	keys := c.sort()

	s := c.heading()
	i := 0
	for _, k := range keys {
		if i%7 == 0 {
			if i > 0 {
				s += c.closeTR()
			}
			s += c.openTR()
		}
		s += c.entry(k)
		i++
	}
	s += c.closeTR()
	return s
}

// New creates a new month calendar based on d, d must be in the form: YYYY-MM-DD.
// D can also be the empty string, then the current date is assumed.
func New(d, subject, secret string) (*Calendar, error) {
	date, now := time.Now(), time.Now()
	if d != "" {
		var err error
		date, err = time.Parse("2006-01-02", d)
		if err != nil {
			return nil, err
		}
	}

	cal := &Calendar{days: make(map[time.Time]AvailMeta), start: date, subject: subject, secret: secret}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	first := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	last = last.Add(-24 * time.Hour)

	// Add the remaining days of the previous month.
	for i := 0; i < int(first.Weekday()); i++ {
		lastMonthDay := first.AddDate(0, 0, -1*(i+1))
		cal.days[lastMonthDay] = AvailMeta{Available: Free, Meta: Prev}

		if lastMonthDay.Before(today) {
			cal.days[lastMonthDay] = AvailMeta{Available: Past, Meta: Prev}
		}
	}

	// Loop from i to lastDay and add the entire month.
	for i := 1; i <= last.Day(); i++ {
		day := time.Date(date.Year(), date.Month(), i, 0, 0, 0, 0, time.UTC)

		cal.days[day] = AvailMeta{Available: Free}

		if day.Before(today) {
			cal.days[day] = AvailMeta{Available: Past}
		}
	}

	// These are dates in the new month.
	j := 1
	for i := int(last.Weekday()) + 1; i < 7; i++ {
		nextMonthDay := last.AddDate(0, 0, j)
		cal.days[nextMonthDay] = AvailMeta{Available: Free, Meta: Next}

		if nextMonthDay.Before(today) {
			cal.days[nextMonthDay] = AvailMeta{Available: Past, Meta: Next}
		}

		j++
	}

	if cur, ok := cal.days[today]; ok {
		cur.Meta = Now
		cal.days[today] = cur
	}

	times := cal.sort()
	if len(times) > 0 {
		cal.begin = times[0]
		cal.end = times[len(times)-1]
	}

	return cal, nil
}
