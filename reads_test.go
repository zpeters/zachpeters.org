package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseRead(t *testing.T) {
	src := "---\ntitle: Slow\nauthor: Derek Sivers\nurl: https://sive.rs/slow\nadded: 2026-10-19\n---\n\nTake your time[^1].\n\n[^1]: Really.\n"
	r, err := ParseRead("slow", []byte(src))
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "slow" || r.Title != "Slow" || r.Author != "Derek Sivers" || r.URL != "https://sive.rs/slow" || r.Year != 0 {
		t.Errorf("read = %+v", r)
	}
	if !r.Added.Equal(time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("added = %v", r.Added)
	}
	if !strings.Contains(string(r.Notes), `id="sn-slow-1"`) {
		t.Errorf("sidenote id not prefixed: %s", r.Notes)
	}
}

func TestParseReadErrors(t *testing.T) {
	for name, src := range map[string]string{
		"no front matter": "Just text.\n",
		"not closed":      "---\ntitle: X\nurl: https://x\n",
		"no url":          "---\ntitle: X\nadded: 2026-10-19\n---\n",
		"bad date":        "---\ntitle: X\nurl: https://x\nadded: someday\n---\n",
	} {
		if _, err := ParseRead("x", []byte(src)); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}

func TestLoadReads(t *testing.T) {
	reads, err := LoadReads("reads/")
	if err != nil {
		t.Fatal(err)
	}
	if len(reads) == 0 {
		t.Fatal("no reads")
	}
}

func TestNewReadTimeline(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	reads := []Read{
		{ID: "slow", Title: "Slow", Year: 2021},
		{ID: "boring", Title: "Choose Boring Technology", Year: 2015},
		{ID: "undated", Title: "Undated"},
	}
	tl := NewReadTimeline(reads, now)
	if len(tl.Rows) != 2 {
		t.Fatalf("got %d rows, want the 2 dated reads", len(tl.Rows))
	}
	if tl.Rows[0].ID != "boring" || tl.Rows[0].Years != "2015" || tl.Rows[0].Line {
		t.Errorf("first row = %+v", tl.Rows[0])
	}
	if tl.StartLabel != "2014" {
		t.Errorf("axis starts %s", tl.StartLabel)
	}
	if NewReadTimeline([]Read{{ID: "undated"}}, now) != nil {
		t.Error("want no timeline when no read has a year")
	}
}
