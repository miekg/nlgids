package ecb

import (
	"bufio"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const ecbURL = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml"

var rateRe = regexp.MustCompile(`currency='GBP' *rate='([.0-9]+)'`)

// RateGBP get the current exchange rate for the Pound from the ECB website.
// The rate returned is GBP:EUR.
func RateGBP() (float64, error) {
	// In good fashion we parse the XML with regular expressions.
	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(ecbURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	rateString := "0.0"
	for scanner.Scan() {
		// The last value of this slice is the rate
		matches := rateRe.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		rateString = matches[len(matches)-1]
		break
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	rate, err := strconv.ParseFloat(rateString, 64)
	if err != nil {
		return 0, nil
	}
	return 1.0 / rate, nil
}
