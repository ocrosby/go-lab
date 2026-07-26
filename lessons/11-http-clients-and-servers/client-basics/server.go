// Package clientbasics demonstrates the smallest useful HTTP client in Go:
// one function that GETs a URL, decodes JSON, and returns a typed result.
package clientbasics

import (
	"encoding/json"
	"net/http"
)

// Weather is the response shape from the mock weather endpoint the tests use.
type Weather struct {
	City     string `json:"city"`
	Forecast string `json:"forecast"`
}

// GetWeather issues an HTTP GET against url and decodes the response body as
// JSON into a *Weather. This is intentionally the simplest possible client —
// no timeout, no context, no retry. Lesson 18 shows the production-ready
// version of these calls.
func GetWeather(url string) (*Weather, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var weather Weather
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, err
	}
	return &weather, nil
}
