package middleware

import (
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	layout := "02-01-2006 15:04" //dd-mm-yyyy hour:minute
	parseTime, err := time.Parse(fmt.Sprintf("\"%s\"", layout), string(b))
	if err != nil {
		return fmt.Errorf("date time format input invalid. expected format: %s", layout)
	}

	c.Time = parseTime
	return nil
}
