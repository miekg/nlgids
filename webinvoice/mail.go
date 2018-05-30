package webinvoice

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"
	"time"

	"github.com/miekg/nlgids/date"
)

// Invoice is a customer invoice form.
type InvoiceMail struct {
	Name    string
	Kenmerk string
	Link    string
}

const templ = `Hallo Ans,

Dit is het reserverings formulier voor "{{.Name}}". Met kenmerk ""{{.Kenmerk}}".

Google Calendar Link:
{{.Link}}

Met vriendelijke groet,
    NLgids mailer
`

func (i *Invoice) MailBody() (*bytes.Buffer, error) {
	t := template.New("Invoice template")
	t, err := t.Parse(templ)
	if err != nil {
		return nil, err
	}

	t1, t2, err := date.TimeToUTC(i.OrigDate, i.OrigTime)
	if err != nil {
		return nil, err
	}

	c := &InvoiceMail{Name: i.FullName, Kenmerk: i.Kenmerk, Link: i.CalendarLink(t1, t2)}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, c); err != nil {
		return nil, err
	}
	return buf, nil
}

func (i *Invoice) MailSubject() string {
	subject := "Formulier (" + i.Kenmerk + "): \"" + i.FullName + "\""
	return subject
}

// CalendarLink returns a Google calendar link that can be used to add the tour
// directly in the calendar.
func (i *Invoice) CalendarLink(t1, t2 time.Time) string {
	dates := "dates=" + url.QueryEscape(fmt.Sprintf("%s/%s", t1.Format(calFmt), t2.Format(calFmt)))
	text := "text=" + url.QueryEscape(fmt.Sprintf("%s - %s", i.FullName, i.Tour))
	loc := "location=London"
	details := "details=" + url.QueryEscape(fmt.Sprintf("Kenmerk: %s", i.Kenmerk))

	link := fmt.Sprintf("%s%s&%s&%s&%s", calLink, dates, text, loc, details)

	return link
}

const (
	calLink = "https://www.google.com/calendar/event?action=TEMPLATE&"
	calFmt  = "20060102T150405Z"
)
