package main

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// ProjectPage is one project page. The projects list shows its Summary
// and links to the whole page at /projects/<ID>.
type ProjectPage struct {
	ID      string
	Name    string
	Years   string // the text in parentheses, e.g. "2015" or "2017-2019"
	Start   int    // first year, for sorting; 0 when unknown
	Summary template.HTML
	Body    template.HTML // the whole file
}

var projectSummary = regexp.MustCompile(`(?s)<p class="summary">\s*(.*?)\s*</p>`)

// LoadProjectPages reads the project pages that match pattern, newest
// first. A page needs a project heading (see projects.go) to count.
// A page without a <p class="summary"> is listed without a summary.
func LoadProjectPages(pattern string) ([]ProjectPage, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var projects []ProjectPage
	for _, path := range paths {
		src, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		text := string(src)

		heading := projectHeading.FindStringSubmatch(text)
		if heading == nil {
			continue
		}
		p := ProjectPage{
			ID:    heading[1],
			Name:  heading[2],
			Years: heading[3],
			Body:  template.HTML(text),
		}
		if parsed := ParseProjects(text); len(parsed) > 0 {
			p.Start = parsed[0].Start
		}
		if s := projectSummary.FindStringSubmatch(text); s != nil {
			p.Summary = template.HTML(s[1])
		}
		projects = append(projects, p)
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Start > projects[j].Start
	})
	return projects, nil
}
