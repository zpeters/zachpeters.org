package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/goccy/go-yaml"
)

// ProjectPage is one project. Each projects/<id>.md file has front matter
// (name, years and a summary) and then the text of the page.
type ProjectPage struct {
	ID      string        // the file name without .md, e.g. "ds"
	Name    string        // e.g. "Dropsearch"
	Years   string        // e.g. "2015" or "2017-2019"
	Start   int           // first year, for sorting; 0 when unknown
	Summary string        // one line, shown on the projects list
	Body    template.HTML // the rendered page, without the heading
}

// projectMeta is the front matter of a projects/ file.
type projectMeta struct {
	Name    string `yaml:"name"`
	Years   string `yaml:"years"`
	Summary string `yaml:"summary"`
}

var projectFence = []byte("---\n")

// ParseProjectPage reads one projects/ file. The front matter must give a
// name and years. A summary is optional.
func ParseProjectPage(id string, src []byte) (ProjectPage, error) {
	if !bytes.HasPrefix(src, projectFence) {
		return ProjectPage{}, fmt.Errorf("%s: no front matter", id)
	}
	meta, body, ok := bytes.Cut(src[len(projectFence):], append([]byte("\n"), projectFence...))
	if !ok {
		return ProjectPage{}, fmt.Errorf("%s: front matter is not closed with ---", id)
	}

	var m projectMeta
	if err := yaml.Unmarshal(meta, &m); err != nil {
		return ProjectPage{}, fmt.Errorf("%s: front matter: %w", id, err)
	}
	if m.Name == "" || m.Years == "" {
		return ProjectPage{}, fmt.Errorf("%s: front matter needs a name and years", id)
	}

	p := ProjectPage{
		ID:      id,
		Name:    m.Name,
		Years:   m.Years,
		Summary: m.Summary,
		Body:    template.HTML(RenderMarkdown(body)),
	}
	if proj, ok := parseProject(id, m.Name, m.Years); ok {
		p.Start = proj.Start
	}
	return p, nil
}

// LoadProjectPages reads the project files that match pattern, newest first.
// A file with years the timeline cannot read is still listed, without a
// place on the timeline.
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
		id := strings.TrimSuffix(filepath.Base(path), ".md")
		p, err := ParseProjectPage(id, src)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Start > projects[j].Start
	})
	return projects, nil
}
