package webinvoice

import (
	"testing"
	"time"
)

func TestKenmerk(t *testing.T) {
	s := Kenmerk(time.Now().UTC())

	t.Logf("%s\n", s)
}
