package tour

import (
	"testing"
)

func TestJsonParse(t *testing.T) {
	tours, err := parse("/home/miek/html/nlgids.london/tours.json")
	if err != nil {
		t.Fatal(err)
	}
	tour := jsonToTour(tours)
	for k, v := range tour {
		t.Logf("%s: %s\n", k, v)
	}
}
