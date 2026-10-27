package site

import (
	"testing"
	"time"
)

var projectsFixture = []ProjectPage{
	{ID: "ok", Name: "Ok - The Git Helper", Years: "Ongoing"},
	{ID: "ds", Name: "Dropsearch", Years: "2024-Current"},
	{ID: "radio", Name: "Radio", Years: "2017-2019"},
	{ID: "os", Name: "Oblique Strategies", Years: "2022"},
	{ID: "bad", Name: "Odd", Years: "someday"},
}

func TestProjectsFromPages(t *testing.T) {
	got := ProjectsFromPages(projectsFixture)
	want := []Project{
		{ID: "ok", Name: "Ok - The Git Helper", Active: true},
		{ID: "ds", Name: "Dropsearch", Start: 2024, Active: true},
		{ID: "radio", Name: "Radio", Start: 2017, End: 2019},
		{ID: "os", Name: "Oblique Strategies", Start: 2022, End: 2022},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d projects: %+v", len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("project %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestNewProjectTimeline(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tl := NewProjectTimeline(ProjectsFromPages(projectsFixture), now)

	if tl.StartLabel != "2016" {
		t.Errorf("axis starts %s, want the year before the first start", tl.StartLabel)
	}
	byID := map[string]ProjectRow{}
	for _, r := range tl.Rows {
		byID[r.ID] = r
	}
	if tl.Rows[0].ID != "ok" {
		t.Errorf("unknown start should sort first, got %s", tl.Rows[0].ID)
	}
	if r := byID["ok"]; !r.Lead || !r.Line || r.X2 != timelineRight || r.Anchor != "start" || r.LabelY >= r.Y {
		t.Errorf("ongoing row = %+v", r)
	}
	if r := byID["ds"]; !r.Line || !r.StartDot || r.EndDot || r.X2 != timelineRight || r.Years != "2024 to now" {
		t.Errorf("active row = %+v", r)
	}
	if r := byID["radio"]; !r.Line || !r.StartDot || !r.EndDot || r.X2 <= r.X1 || r.Years != "2017 to 2019" {
		t.Errorf("range row = %+v", r)
	}
	if r := byID["os"]; r.Line || !r.StartDot || r.X1 != r.X2 || r.Years != "2022" {
		t.Errorf("single-year row = %+v", r)
	}
	if NewProjectTimeline(nil, now) != nil {
		t.Error("want nil timeline with no projects")
	}
}
