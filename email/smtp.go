package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

const (
	server       = "localhost:25"
	nlgidsPrefix = "[NLgids] "
)

var (
	//to = []string{"ans@nlgids.london", "miek@miek.nl"}
	to   = []string{"miek@miek.nl"}
	from = "nlgids@x64.atoom.net"
)

func NewContact(subject string, body *bytes.Buffer) *Message {
	m := NewMessage(nlgidsPrefix+subject, body.String())
	m.From = from
	m.To = to
	return m
}

func NewBooking(subject string, body *bytes.Buffer) *Message {
	m := NewMessage(nlgidsPrefix+subject, body.String())
	m.From = from
	m.To = to
	return m
}

func NewInvoice(subject string, body *bytes.Buffer, pdfName string, pdf []byte) *Message {
	m := NewMessage(nlgidsPrefix+subject, body.String())
	m.From = from
	m.To = to
	m.AttachBytes(pdfName, pdf)
	return m
}

// Do sends an email message.
func (m *Message) Do() error { return Send(server, nil, m) }

// Everything below:

// Copyright 2012 Santiago Corredoira
// Distributed under a BSD-like license.

type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

type Message struct {
	From            string
	To              []string
	Cc              []string
	Bcc             []string
	ReplyTo         string
	Subject         string
	Body            string
	BodyContentType string
	Attachments     map[string]*Attachment
}

func (m *Message) attach(file string, inline bool) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	m.attachBytes(file, data, inline)
	return nil
}

func (m *Message) attachBytes(file string, data []byte, inline bool) {
	_, filename := filepath.Split(file)

	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}
}

func (m *Message) AttachBytes(file string, data []byte) {
	m.attachBytes(file, data, false)
}

func (m *Message) InlineBytes(file string, data []byte) {
	m.attachBytes(file, data, true)
}

func (m *Message) Attach(file string) error {
	return m.attach(file, false)
}

func (m *Message) Inline(file string) error {
	return m.attach(file, true)
}

func newMessage(subject string, body string, bodyContentType string) *Message {
	m := &Message{Subject: subject, Body: body, BodyContentType: bodyContentType}

	m.Attachments = make(map[string]*Attachment)

	return m
}

// NewMessage returns a new Message that can compose an email with attachments
func NewMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/plain")
}

// NewMessage returns a new Message that can compose an HTML email with attachments
func NewHTMLMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/html")
}

// ToList returns all the recipients of the email
func (m *Message) Tolist() []string {
	tolist := m.To

	for _, cc := range m.Cc {
		tolist = append(tolist, cc)
	}

	for _, bcc := range m.Bcc {
		tolist = append(tolist, bcc)
	}

	return tolist
}

// Bytes returns the mail data
func (m *Message) Bytes() []byte {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("From: " + m.From + "\r\n")

	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC822) + "\r\n")

	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.Cc) > 0 {
		buf.WriteString("Cc: " + strings.Join(m.Cc, ",") + "\r\n")
	}

	buf.WriteString("Subject: " + m.Subject + "\r\n")

	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	}

	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.BodyContentType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")

			if attachment.Inline {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				buf.Write(attachment.Data)
			} else {
				buf.WriteString("Content-Type: application/octet-stream\r\n")
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString("Content-Disposition: attachment; filename=\"" + attachment.Filename + "\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
				base64.StdEncoding.Encode(b, attachment.Data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

func Send(addr string, auth smtp.Auth, m *Message) error {
	return smtp.SendMail(addr, auth, m.From, m.Tolist(), m.Bytes())
}
