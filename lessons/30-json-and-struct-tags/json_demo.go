// Package jsondemo demonstrates encoding/json — Marshal/Unmarshal, the
// four common struct-tag options, safe body decoding for HTTP handlers,
// and a custom MarshalJSON for a domain type.
package jsondemo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// User illustrates the four common struct-tag options.
type User struct {
	ID    string `json:"id"`                // renamed
	Email string `json:"email"`             // renamed
	Age   int    `json:"age,omitempty"`     // rename + omit when zero
	Notes string `json:"-"`                 // never serialized
	Money Money  `json:"balance,omitempty"` // uses custom MarshalJSON below
}

// Money is an int-cents amount that renders in JSON as a dollar-cents
// string like "12.34". Demonstrates a custom MarshalJSON/UnmarshalJSON
// pair — the domain type's wire format differs from its Go shape.
type Money int64

// MarshalJSON renders Money as a JSON string of dollars.cents.
func (m Money) MarshalJSON() ([]byte, error) {
	dollars := int64(m) / 100
	cents := int64(m) % 100
	if cents < 0 {
		cents = -cents
	}
	return []byte(fmt.Sprintf("%q", fmt.Sprintf("%d.%02d", dollars, cents))), nil
}

// UnmarshalJSON parses "12.34" back into Money (1234 cents).
func (m *Money) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parts := strings.SplitN(s, ".", 2)
	dollars, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Money.UnmarshalJSON: dollars: %w", err)
	}
	var cents int64
	if len(parts) == 2 {
		cents, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("Money.UnmarshalJSON: cents: %w", err)
		}
	}
	*m = Money(dollars*100 + cents)
	return nil
}

// DecodeStrict is the safe body-decoding pattern from lesson 16, in
// standalone form. Caller passes any io.Reader (an http.Request.Body,
// a strings.Reader in a test, etc.). Streaming decoder catches oversized
// bodies as they arrive; DisallowUnknownFields catches caller typos;
// the More() check catches concatenated-JSON inputs.
//
// maxBytes bounds the number of bytes read. Every request handler
// should cap this — an unbounded body is an OOM vector.
func DecodeStrict(r io.Reader, maxBytes int64, dst any) error {
	limited := io.LimitReader(r, maxBytes+1)
	dec := json.NewDecoder(limited)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	if dec.More() {
		return errors.New("body must contain exactly one JSON value")
	}
	return nil
}
