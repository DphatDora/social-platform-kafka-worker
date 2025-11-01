package util

import (
	"fmt"
	"time"
)

func FormatMonthYear(t time.Time) string {
	return fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
}
