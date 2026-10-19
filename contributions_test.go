package main

import (
	"strings"
	"testing"
	"time"
)

const calendarFixture = `<table><tbody>
<tr>
<td data-date="2025-10-19" id="contribution-day-component-0-0" class="ContributionCalendar-day"></td>
<td data-date="2025-10-26" id="contribution-day-component-0-1" class="ContributionCalendar-day"></td>
</tr><tr>
<td data-date="2025-10-20" id="contribution-day-component-1-0" class="ContributionCalendar-day"></td>
<td data-date="2025-10-27" id="contribution-day-component-1-1" class="ContributionCalendar-day"></td>
</tr></tbody></table>
<tool-tip for="contribution-day-component-0-0" class="sr-only">No contributions on October 19th.</tool-tip>
<tool-tip for="contribution-day-component-1-0" class="sr-only">1 contribution on October 20th.</tool-tip>
<tool-tip for="contribution-day-component-0-1" class="sr-only">1,200 contributions on October 26th.</tool-tip>
<tool-tip for="contribution-day-component-1-1" class="sr-only">3 contributions on October 27th.</tool-tip>
`

func TestParseContributions(t *testing.T) {
	weeks, from, to, err := ParseContributions(strings.NewReader(calendarFixture))
	if err != nil {
		t.Fatal(err)
	}
	if len(weeks) != 2 || weeks[0] != 1 || weeks[1] != 1203 {
		t.Errorf("weeks = %v, want [1 1203]", weeks)
	}
	if got := from.Format(time.DateOnly); got != "2025-10-19" {
		t.Errorf("from = %s", got)
	}
	if got := to.Format(time.DateOnly); got != "2025-10-27" {
		t.Errorf("to = %s", got)
	}
}

func TestParseContributionsEmpty(t *testing.T) {
	if _, _, _, err := ParseContributions(strings.NewReader("<html></html>")); err == nil {
		t.Error("want error for a page with no calendar")
	}
}

func TestNewSparkline(t *testing.T) {
	s := NewSparkline([]int{0, 10, 5}, time.Time{}, time.Time{})
	if s.Points != "0.0,36.0 120.0,0.0 240.0,18.0" {
		t.Errorf("points = %q", s.Points)
	}
	if s.Max != 10 || s.MaxX != 120 || s.MaxY != 0 {
		t.Errorf("max = %d at (%v,%v)", s.Max, s.MaxX, s.MaxY)
	}
	if s.Last != 5 || s.LastX != 240 || s.LastY != 18 {
		t.Errorf("last = %d at (%v,%v)", s.Last, s.LastX, s.LastY)
	}
	if s.Total != "15" {
		t.Errorf("total = %q", s.Total)
	}
}

func TestThousands(t *testing.T) {
	for n, want := range map[int]string{0: "0", 999: "999", 4142: "4,142", 1234567: "1,234,567"} {
		if got := thousands(n); got != want {
			t.Errorf("thousands(%d) = %q, want %q", n, got, want)
		}
	}
}
