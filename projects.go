package main

import (
	"fmt"
	"html/template"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Project is one project on the timeline. Its years come from the page's
// front matter (see project_pages.go).
type Project struct {
	ID     string
	Name   string
	Start  int  // 0 when unknown ("Ongoing")
	End    int  // last year; 0 when Active
	Active bool // still going
}

var projectYears = regexp.MustCompile(`(?i)^(\d{4})(?:\s*[-–]\s*(\d{4}|current|present|now))?$`)

// parseProject makes a Project from a page's name and years. The years are
// "2022", "2017-2019", "2023-Current" or "Ongoing". Other years are not
// on the timeline, and parseProject returns false for them.
func parseProject(id, name, years string) (Project, bool) {
	p := Project{ID: id, Name: name}
	years = strings.TrimSpace(years)
	if strings.EqualFold(years, "ongoing") {
		p.Active = true
		return p, true
	}
	y := projectYears.FindStringSubmatch(years)
	if y == nil {
		return Project{}, false
	}
	p.Start, _ = strconv.Atoi(y[1])
	switch end, err := strconv.Atoi(y[2]); {
	case y[2] == "":
		p.End = p.Start
	case err == nil:
		p.End = end
	default:
		p.Active = true
	}
	return p, true
}

// ProjectsFromPages returns the timeline entries for the pages whose years
// the timeline can read.
func ProjectsFromPages(pages []ProjectPage) []Project {
	var projects []Project
	for _, pg := range pages {
		if p, ok := parseProject(pg.ID, pg.Name, pg.Years); ok {
			projects = append(projects, p)
		}
	}
	return projects
}

// LoadProjectTimeline builds the timeline from the project files that match
// pattern, e.g. "projects/*.md".
func LoadProjectTimeline(pattern string, now time.Time) (*ProjectTimeline, error) {
	pages, err := LoadProjectPages(pattern)
	if err != nil {
		return nil, err
	}
	return NewProjectTimeline(ProjectsFromPages(pages), now), nil
}

// ProjectTimeline lays projects out one per row on a range-frame year axis
// that runs to now, in an SVG TimelineWidth wide and Height tall.
type ProjectTimeline struct {
	Title        string        // the SVG title
	Caption      template.HTML // the margin note
	Height       float64
	AxisY        float64
	TickY        float64 // bottom of the year ticks under the axis
	StartX, EndX float64
	StartLabel   string
	YearTicks    []float64
	Rows         []ProjectRow
}

// ProjectRow is one project: a dot for a single year, a line with a dot at
// each end for a range, or a dot and a line to now for an active project.
// Lead marks an active project whose start is before the axis.
type ProjectRow struct {
	ID, Name, Years  string
	Y                float64
	X1, X2           float64
	Line             bool
	StartDot, EndDot bool
	Lead             bool
	LabelX, LabelY   float64
	Anchor           string
}

const (
	projectRowGap   = 26.0
	projectTop      = 18.0
	projectCharWide = 6.6 // rough label width per character at 13px
)

func NewProjectTimeline(projects []Project, now time.Time) *ProjectTimeline {
	if len(projects) == 0 {
		return nil
	}
	sorted := append([]Project(nil), projects...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Start < sorted[j].Start })

	// The axis starts a year before the earliest known start
	first := now.Year()
	for _, p := range sorted {
		if p.Start != 0 {
			first = min(first, p.Start)
		}
	}
	axisStart := float64(first - 1)
	nowYear := float64(now.Year()) + float64(now.YearDay()-1)/365
	x := func(year float64) float64 {
		return timelineLeft + (timelineRight-timelineLeft)*(year-axisStart)/(nowYear-axisStart)
	}

	t := &ProjectTimeline{
		Title:      "Projects by year",
		Caption:    "Projects by year. A line that runs to <em>now</em> is still active; a dot alone is a single year. Each name links to its section below.",
		StartX:     timelineLeft,
		EndX:       timelineRight,
		StartLabel: strconv.Itoa(first - 1),
		AxisY:      projectTop + float64(len(sorted))*projectRowGap,
	}
	t.TickY = t.AxisY + 4
	t.Height = t.AxisY + 24
	for y := first; float64(y) < nowYear; y++ {
		t.YearTicks = append(t.YearTicks, x(float64(y)))
	}

	for i, p := range sorted {
		r := ProjectRow{ID: p.ID, Name: p.Name, Y: projectTop + float64(i)*projectRowGap}
		switch {
		case p.Active && p.Start == 0:
			r.X1, r.X2, r.Line, r.Lead = timelineLeft, timelineRight, true, true
			r.Years = "ongoing"
		case p.Active:
			r.X1, r.X2, r.Line, r.StartDot = x(float64(p.Start)), timelineRight, true, true
			r.Years = fmt.Sprintf("%d to now", p.Start)
		case p.End > p.Start:
			r.X1, r.X2 = x(float64(p.Start)), min(x(float64(p.End+1)), timelineRight)
			r.Line, r.StartDot, r.EndDot = true, true, true
			r.Years = fmt.Sprintf("%d to %d", p.Start, p.End)
		default:
			r.X1 = x(float64(p.Start) + 0.5)
			r.X2 = r.X1
			r.StartDot = true
			r.Years = strconv.Itoa(p.Start)
		}

		// Label to the left of the mark, else to the right, else above it
		w := float64(len(p.Name)) * projectCharWide
		switch {
		case r.X1-8-w >= 0:
			r.LabelX, r.LabelY, r.Anchor = r.X1-8, r.Y+4, "end"
		case r.X2+8+w <= TimelineWidth:
			r.LabelX, r.LabelY, r.Anchor = r.X2+8, r.Y+4, "start"
		default:
			r.LabelX, r.LabelY, r.Anchor = r.X1, r.Y-7, "start"
		}
		t.Rows = append(t.Rows, r)
	}
	return t
}
