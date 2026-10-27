package site

import (
	"html"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/truncate"
	"github.com/microcosm-cc/bluemonday"
)

var sentenceEnd = regexp.MustCompile(`[.!?](\s|$)`)

// Excerpt returns the first sentence of a markdown post as plain text.
func Excerpt(markdown string) string {
	text := bluemonday.StrictPolicy().Sanitize(string(RenderMarkdown([]byte(markdown))))
	text = strings.Join(strings.Fields(html.UnescapeString(text)), " ")
	if loc := sentenceEnd.FindStringIndex(text); loc != nil {
		return text[:loc[0]+1]
	}
	return truncate.Truncate(text, 200, "…", truncate.PositionEnd)
}

// PostTimeline places blog posts on a range-frame time axis, from the start
// of the first post's year to now, in an SVG TimelineWidth wide.
type PostTimeline struct {
	Marks      []TimelineMark
	StartX     float64
	EndX       float64
	StartLabel string
}

// TimelineMark is one post: a dot at X with its title above, on one of two
// label rows so neighbors don't collide.
type TimelineMark struct {
	X      float64
	LabelY float64
	Leader float64 // top of the leader line from the label down to the dot
	Anchor string  // SVG text-anchor: labels open toward the middle
	Title  string
	URL    string
}

const (
	TimelineWidth = 560.0
	timelineLeft  = 20.0
	timelineRight = TimelineWidth - 20
)

// Label rows sit above the post line at y = 56
var timelineLabelRows = []float64{22, 40}

func NewPostTimeline(posts []BlogPost, now time.Time) *PostTimeline {
	if len(posts) == 0 {
		return nil
	}
	sorted := append([]BlogPost(nil), posts...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Date.Before(sorted[j].Date) })

	start := time.Date(sorted[0].Date.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	span := float64(now.Sub(start))
	t := &PostTimeline{
		StartX:     timelineLeft,
		EndX:       timelineRight,
		StartLabel: start.Format("2006"),
	}
	for i, p := range sorted {
		x := timelineLeft + (timelineRight-timelineLeft)*float64(p.Date.Sub(start))/span
		anchor := "start"
		if x > TimelineWidth/2 {
			anchor = "end"
		}
		t.Marks = append(t.Marks, TimelineMark{
			X:      x,
			LabelY: timelineLabelRows[i%len(timelineLabelRows)],
			Leader: timelineLabelRows[i%len(timelineLabelRows)] + 5,
			Anchor: anchor,
			Title:  p.Title,
			URL:    p.URL,
		})
	}
	return t
}
