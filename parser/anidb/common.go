package anidb

import (
	"errors"
	"strings"
	"time"
)

// Parses a raw date from AniDB anime page. Returns zero for empty string.
func parseDate(s string) (time.Time, error) {
	var zero time.Time

	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return zero, errors.New("date is empty")
	}

	var format string
	if strings.Contains(s, "-") {
		format = "2006-01-02"
	} else {
		format = "02.01.2006"
	}

	return time.Parse(format, s)
}
