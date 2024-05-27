package pkg

import (
	"database/sql/driver"
	"strconv"
	"strings"
	"time"
)

// Handling Floats

type Float64String struct {
	float64
}

func (f *Float64String) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = strings.Trim(str, "\"") // Remove quotation marks

	if str == "" || str == "null" {
		f.float64 = 0
		return nil
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	f.float64 = val
	return nil
}

func (f Float64String) String() string {
	return strconv.FormatFloat(f.float64, 'f', -1, 64)
}

func (f Float64String) Value() (driver.Value, error) {
	return f.float64, nil
}

func NewFloat64String(f float64) Float64String {
	return Float64String{f}
}

// Handling DateTimes

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02T15:04:05.000"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(ctLayout, s)
	if err != nil {
		// If parsing fails with the custom layout, try the RFC 3339 layout
		ct.Time, err = time.Parse(time.RFC3339, s)
	}
	return
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}
