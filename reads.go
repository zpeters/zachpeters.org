package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
)

// Read is an article by someone else that I recommend. Each one is a
// markdown file in reads/: YAML front matter, then my notes on it.
type Read struct {
	ID      string // file name without .md; the anchor on the Reads page
	Title   string
	Author  string
	Year    int    // when it was written; 0 when unknown
	URL     string // the original
	Archive string // a saved copy, in case the original goes away
	Added   time.Time
	Notes   template.HTML
}

type readMeta struct {
	Title   string `yaml:"title"`
	Author  string `yaml:"author"`
	Year    int    `yaml:"year"`
	URL     string `yaml:"url"`
	Archive string `yaml:"archive"`
	Added   string `yaml:"added"`
}

var frontMatterFence = []byte("---\n")

// ParseRead reads one reads/ file. The front matter must give a title, a
// url and the date it was added (YYYY-MM-DD).
func ParseRead(id string, src []byte) (Read, error) {
	src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(src, frontMatterFence) {
		return Read{}, fmt.Errorf("no front matter")
	}
	meta, body, ok := bytes.Cut(src[len(frontMatterFence):], append([]byte("\n"), frontMatterFence...))
	if !ok {
		return Read{}, fmt.Errorf("front matter is not closed with ---")
	}

	var m readMeta
	if err := yaml.Unmarshal(meta, &m); err != nil {
		return Read{}, fmt.Errorf("front matter: %w", err)
	}
	if m.Title == "" || m.URL == "" {
		return Read{}, fmt.Errorf("front matter needs a title and a url")
	}
	added, err := time.Parse("2006-01-02", m.Added)
	if err != nil {
		return Read{}, fmt.Errorf("added: %w", err)
	}
	return Read{
		ID:      id,
		Title:   m.Title,
		Author:  m.Author,
		Year:    m.Year,
		URL:     m.URL,
		Archive: m.Archive,
		Added:   added,
		// Prefix the sidenote ids, so notes in two reads on one page do not clash
		Notes: template.HTML(renderMarkdown(bytes.TrimSpace(body), "sn-"+id+"-")),
	}, nil
}

// LoadReads reads every reads/*.md file, newest added first.
func LoadReads(dir string) ([]Read, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return nil, err
	}
	var reads []Read
	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Skipping read %s: %v", file, err)
			continue
		}
		r, err := ParseRead(strings.TrimSuffix(filepath.Base(file), ".md"), src)
		if err != nil {
			log.Printf("Skipping read %s: %v", file, err)
			continue
		}
		reads = append(reads, r)
	}
	sort.SliceStable(reads, func(i, j int) bool { return reads[i].Added.After(reads[j].Added) })
	return reads, nil
}

// NewReadTimeline lays out the reads by the year they were written, one per
// row, on the same chart as the project pages. Reads without a year are
// left off.
func NewReadTimeline(reads []Read, now time.Time) *ProjectTimeline {
	var dated []Project
	for _, r := range reads {
		if r.Year != 0 {
			dated = append(dated, Project{ID: r.ID, Name: r.Title, Start: r.Year, End: r.Year})
		}
	}
	t := NewProjectTimeline(dated, now)
	if t != nil {
		t.Title = "Reads by the year they were written"
		t.Caption = "Reads by the year they were written. Each title links to my notes on it below."
	}
	return t
}
