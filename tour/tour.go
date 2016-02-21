package tour

import (
	"encoding/json"
	"io/ioutil"
)

const nonExists = "<niet bestaand>"

type Definition struct {
	ID   string
	Naam string
	Name string
}

type Tours map[string][]Definition

func parse(file string) (Tours, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var tours Tours
	if err := json.Unmarshal(f, &tours); err != nil {
		return nil, err
	}
	return tours, nil
}

// jsonToTour molds the Tours json into the Tour var we want:
// cycling/custom: Dutch name
func jsonToTour(t Tours) map[string]string {
	m := map[string]string{}
	for typ := range t {
		for _, tour := range t[typ] {
			key := typ + "/" + tour.ID
			m[key] = tour.Naam
		}
	}
	return m
}

func NewTour(file string) map[string]string {
	t, _ := parse(file)
	m := jsonToTour(t)
	return m
}

// Tour holds the conversion from tour ID to the Dutch name.
var Tour = map[string]string{}

// TourOrNonExists returns either the Dutch name of the tour keyed
// by key or the string "<niet bestaand>".
func NameOrNonExists(key, file string) string {
	if Tour == nil {
		Tour = NewTour(file)
	}

	if v, ok := Tour[key]; ok {
		return v
	}
	return nonExists
}
