package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/html"
)

// Sparkline is a year of weekly GitHub contribution counts, laid out for an
// SVG of SparkWidth x SparkHeight.
type Sparkline struct {
	Points       string // polyline points
	MaxX, MaxY   float64
	LastX, LastY float64
	Max, Last    int
	Total        string // with thousands separators
	From, To     time.Time
}

const (
	SparkWidth  = 560.0
	SparkHeight = 60.0
)

// StartContributions loads the user's contribution calendar now and then
// every interval, and returns a func for the latest sparkline (nil until
// the first load succeeds).
func StartContributions(user string, interval time.Duration) func() *Sparkline {
	var cur atomic.Pointer[Sparkline]
	client := &http.Client{Timeout: 15 * time.Second}
	go func() {
		for {
			s, err := loadSparkline(context.Background(), client, user)
			if err != nil {
				log.Printf("Could not load GitHub contributions: %v", err)
			} else {
				cur.Store(s)
			}
			time.Sleep(interval)
		}
	}()
	return cur.Load
}

func loadSparkline(ctx context.Context, client *http.Client, user string) (*Sparkline, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://github.com/users/"+user+"/contributions", nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %s", resp.Status)
	}
	weeks, from, to, err := ParseContributions(resp.Body)
	if err != nil {
		return nil, err
	}
	return NewSparkline(weeks, from, to), nil
}

// ParseContributions reads GitHub's contribution calendar HTML and returns
// the total per week (calendar column), oldest first, and the first and
// last dates. Each day is a <td id=... data-date=...>; its count is in the
// <tool-tip for=id> text, e.g. "No contributions on ..." or "12 contributions on ...".
func ParseContributions(r io.Reader) ([]int, time.Time, time.Time, error) {
	const idPrefix = "contribution-day-component-"
	var (
		weekOf   = map[string]int{} // td id -> week column
		counts   = map[string]int{} // td id -> count
		weeks    []int
		from, to time.Time
		tipFor   string
	)

	z := html.NewTokenizer(r)
	for {
		switch z.Next() {
		case html.ErrorToken:
			if z.Err() != io.EOF {
				return nil, from, to, z.Err()
			}
			if len(weekOf) == 0 {
				return nil, from, to, fmt.Errorf("no contribution days found")
			}
			for id, w := range weekOf {
				for len(weeks) <= w {
					weeks = append(weeks, 0)
				}
				weeks[w] += counts[id]
			}
			return weeks, from, to, nil
		case html.StartTagToken:
			t := z.Token()
			switch t.Data {
			case "td":
				id, date := attr(t, "id"), attr(t, "data-date")
				if !strings.HasPrefix(id, idPrefix) || date == "" {
					continue
				}
				// id is contribution-day-component-<weekday>-<week>
				parts := strings.Split(strings.TrimPrefix(id, idPrefix), "-")
				if len(parts) != 2 {
					continue
				}
				w, err := strconv.Atoi(parts[1])
				if err != nil {
					continue
				}
				weekOf[id] = w
				if d, err := time.Parse("2006-01-02", date); err == nil {
					if from.IsZero() || d.Before(from) {
						from = d
					}
					if d.After(to) {
						to = d
					}
				}
			case "tool-tip":
				tipFor = attr(t, "for")
			}
		case html.TextToken:
			if tipFor != "" {
				counts[tipFor] = parseCount(string(z.Text()))
				tipFor = ""
			}
		}
	}
}

func attr(t html.Token, name string) string {
	for _, a := range t.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// parseCount reads the leading number of a tooltip like "1,234 contributions on ...".
func parseCount(s string) int {
	n, _ := strconv.Atoi(strings.ReplaceAll(strings.Fields(s + " 0")[0], ",", ""))
	return n
}

// NewSparkline lays the weekly counts out left to right, highest at the top.
func NewSparkline(weeks []int, from, to time.Time) *Sparkline {
	s := &Sparkline{From: from, To: to}
	total := 0
	for _, v := range weeks {
		total += v
		s.Max = max(s.Max, v)
	}
	s.Total = thousands(total)

	pts := make([]string, len(weeks))
	for i, v := range weeks {
		x := SparkWidth * float64(i) / float64(max(len(weeks)-1, 1))
		y := SparkHeight
		if s.Max > 0 {
			y = SparkHeight * (1 - float64(v)/float64(s.Max))
		}
		pts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
		if v == s.Max && s.MaxX == 0 && s.MaxY == 0 {
			s.MaxX, s.MaxY = x, y
		}
		if i == len(weeks)-1 {
			s.LastX, s.LastY, s.Last = x, y, v
		}
	}
	s.Points = strings.Join(pts, " ")
	return s
}

func thousands(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}
