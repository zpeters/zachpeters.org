package site

import (
	"testing"
	"time"
)

func TestExcerpt(t *testing.T) {
	for md, want := range map[string]string{
		`<span class="newthought">There is something</span> about it.  Second one.`: "There is something about it.",
		"No full stop here":                       "No full stop here",
		"Uses \"quotes\" and *style*! Then more.": "Uses “quotes” and style!",
		"Version 1.2 is out. Next.":               "Version 1.2 is out.",
	} {
		if got := Excerpt(md); got != want {
			t.Errorf("Excerpt(%q) = %q, want %q", md, got, want)
		}
	}
}

func TestNewPostTimeline(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	posts := []BlogPost{
		{Title: "Newer", URL: "newer", Date: time.Date(2025, 7, 2, 12, 0, 0, 0, time.UTC)},
		{Title: "Older", URL: "older", Date: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	tl := NewPostTimeline(posts, now)
	if tl.StartLabel != "2024" {
		t.Errorf("start label = %q", tl.StartLabel)
	}
	if len(tl.Marks) != 2 || tl.Marks[0].Title != "Older" || tl.Marks[1].Title != "Newer" {
		t.Fatalf("marks not oldest first: %+v", tl.Marks)
	}
	if m := tl.Marks[0]; m.X != timelineLeft || m.Anchor != "start" || m.LabelY != timelineLabelRows[0] {
		t.Errorf("older mark = %+v", m)
	}
	if m := tl.Marks[1]; m.X <= TimelineWidth/2 || m.Anchor != "end" || m.LabelY != timelineLabelRows[1] {
		t.Errorf("newer mark = %+v", m)
	}
	if NewPostTimeline(nil, now) != nil {
		t.Error("want nil timeline with no posts")
	}
}
