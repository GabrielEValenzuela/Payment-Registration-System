package models

import (
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const customTimeFormat = "2006-01-02T15:04:05.999999"

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
