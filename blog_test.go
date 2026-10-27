package main

import (
	"testing"
	"time"
)

func TestParsePost(t *testing.T) {
	title, date, body, err := parsePost([]byte("---\ntitle: Tools_and - Notes\ndate: \"2026-10-18T21:30:00-0500\"\n---\n\nHello.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if title != "Tools_and - Notes" || body != "Hello." {
		t.Errorf("title = %q, body = %q", title, body)
	}
	if !date.Equal(time.Date(2026, 10, 18, 21, 30, 0, 0, time.FixedZone("CDT", -5*3600))) {
		t.Errorf("date = %v", date)
	}
	for name, src := range map[string]string{
		"no front matter": "Hello.",
		"not closed":      "---\ntitle: X\ndate: \"2026-10-18T21:30:00-0500\"\nHello.",
		"no title":        "---\ndate: \"2026-10-18T21:30:00-0500\"\n---\nHello.",

		"no date":   "---\ntitle: X\n---\nHello.",
		"date only": "---\ntitle: X\ndate: \"2026-10-18\"\n---\nHello.",
	} {
		if _, _, _, err := parsePost([]byte(src)); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
}

func TestParseTimestamp(t *testing.T) {
	for _, s := range []string{"2026-10-18T21:30:00-0500", "2026-10-18T21:30:00-05:00"} {
		got, err := parseTimestamp(s)
		if err != nil {
			t.Errorf("%s: %v", s, err)
			continue
		}
		if !got.Equal(time.Date(2026, 10, 19, 2, 30, 0, 0, time.UTC)) {
			t.Errorf("%s = %v", s, got)
		}
	}
	if _, err := parseTimestamp("2026-10-18"); err == nil {
		t.Error("date only: want error")
	}
}
