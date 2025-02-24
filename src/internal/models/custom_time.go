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
	"time"
)

// CustomTime represents a custom time type that allows parsing time values in a specific format.
type CustomTime struct {
	time.Time
}

const customTimeFormat = "2006-01-02T15:04:05.999999"

// MarshalJSON converts the CustomTime value to a JSON string.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// Eliminar las comillas del valor JSON
	str := string(b)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	// Parsear el tiempo usando el formato personalizado
	parsedTime, err := time.Parse(customTimeFormat, str)
	if err != nil {
		return fmt.Errorf("invalid time format: %v", err)
	}

	ct.Time = parsedTime
	return nil
}
