package date

import (
	"testing"
	"time"
)

func TestTimeToUTC(t *testing.T) {
	t1, _, err := TimeToUTC("2018/07/06", "9.45")

	if err != nil {
		t.Errorf("Failed to convert times: %s", err)
	}

	if t1.Format(time.RFC3339) != "2018-07-06T08:45:00Z" {
		t.Errorf("Failed to parse date")
	}
}
