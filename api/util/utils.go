package util

import (
	"context"
	"time"

	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.Get()

// ParseTime_RFC3339Nano parses timestamps in format "2024-03-18T10:00:00.0000000Z" to go time
func ParseTime_RFC3339Nano(val string) time.Time {
	time1, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		logger.WithContext(context.Background()).Warnf("failed to parse '%s' to time", val)
	}
	return time1
}

// ToFormatted formats the given go time to string in format "2024-03-18T10:00:00.0000000Z"
func ToFormatted(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

// ToFormattedUTC converts the given go time to UTC, then formats it to string in format "2024-03-18T10:00:00.0000000Z"
func ToFormattedUTC(t time.Time) string {
	return ToFormatted(t.UTC())
}
