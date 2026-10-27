package site

import (
	"fmt"
	"time"
)

// timestampLayouts are the forms a header date can take. The UTC offset may
// be written with or without a colon.
var timestampLayouts = []string{"2006-01-02T15:04:05-0700", time.RFC3339}

// parseTimestamp reads the date in a content file's header, for example
// "2026-10-22T09:40:00-0500". A date with no time is an error, so every
// file records when it was written.
func parseTimestamp(s string) (time.Time, error) {
	for _, layout := range timestampLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("date %q is not a timestamp like 2026-10-22T09:40:00-0500", s)
}
