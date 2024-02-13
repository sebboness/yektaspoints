package util

import (
	"context"
	"time"

	"github.com/sebboness/yektaspoints/util/log"
)

var logger = log.Get()

func ParseTime_RFC3339Nano(val string) time.Time {
	time1, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		logger.WithContext(context.Background()).Warnf("failed to parse '%s' to time", val)
	}
	return time1
}

func ToFormattedUTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
