/*
 * Payment Registration System - Custom Time
 * -----------------------------------------
 * This file defines a custom time type that allows parsing time values
 * in a specific format for JSON unmarshalling.
 *
 * Created: Oct. 19, 2024
 * License: GNU General Public License v3.0
 */

package models

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime struct to handle time parsing correctly
type CustomTime struct {
	time.Time
}

const customTimeFormat = time.RFC3339 // Use standard RFC3339 format

// MarshalJSON converts CustomTime to JSON as an RFC3339 string.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte(`null`), nil // Properly handle zero values
	}
	return []byte(`"` + ct.Time.Format(customTimeFormat) + `"`), nil
}

// UnmarshalJSON parses an RFC3339 time string into CustomTime.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Remove quotes from the JSON string
	str := strings.Trim(string(b), `"`)

	// If empty, return zero time
	if str == "null" || str == "" {
		ct.Time = time.Time{}
		return nil
	}

	// Parse time using RFC3339
	parsedTime, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}

	ct.Time = parsedTime
	return nil
}
