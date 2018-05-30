package date

import "time"

// TimeToUTC converts date (YYYY/MM/DD) and time (HH.MM) which are in the London
// time zone to UTC. Is returns two UTC timestamp where the second one is 4 hours later
// than the first.
func TimeToUTC(date string, t string) (time.Time, time.Time, error) {
	form := "2006/01/02 15.04"
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	t1, err := time.ParseInLocation(form, date+" "+t, london)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	t2 := t1.Add(fourHours)

	return t1.UTC(), t2.UTC(), nil
}

const fourHours = time.Duration(4 * time.Hour)
