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

	"github.com/aquilax/truncate"
	"github.com/goccy/go-yaml"
	"github.com/gorilla/feeds"
	"github.com/microcosm-cc/bluemonday"
)

type BlogPost struct {
	Date          time.Time
	FormattedDate string
	Title         string
	URL           string // the file name without .md, e.g. "coding-for-one"
	Content       string
	HTML          template.HTML // Content rendered to HTML, once, at load
	Excerpt       string        // first sentence, plain text
}

func LoadBlogPosts(postsDir string) ([]BlogPost, error) {
	var posts []BlogPost

	files, err := filepath.Glob(filepath.Join(postsDir, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		base := filepath.Base(file)
		src, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Skipping file %s: %v", base, err)
			continue
		}
		title, date, content, err := parsePost(src)
		if err != nil {
			log.Printf("Skipping file %s: %v", base, err)
			continue
		}

		posts = append(posts, BlogPost{
			Date:          date,
			FormattedDate: FormatDate(date),
			Title:         title,
			URL:           strings.TrimSuffix(base, ".md"),
			Content:       content,
			HTML:          template.HTML(RenderMarkdown([]byte(content))),
			Excerpt:       Excerpt(content),
		})
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func GenerateRssFeed(posts []BlogPost) (string, error) {
	feed := &feeds.Feed{
		Title:       "Zach Peters Website",
		Link:        &feeds.Link{Href: "https://zachpeters.org"},
		Description: "A personal website",
		Author: &feeds.Author{
			Name:  "Zach Peters",
			Email: "zpeters@gmail.com",
		},
		Created: time.Now(),
	}

	for _, post := range posts {
		p := bluemonday.StripTagsPolicy()
		// Strip tags from the rendered HTML, so markdown syntax does not show
		strippedHtml := p.Sanitize(string(post.HTML))
		feed.Items = append(feed.Items, &feeds.Item{
			Title: post.Title,
			Link:  &feeds.Link{Href: "https://zachpeters.org/blog/" + post.URL},
			Author: &feeds.Author{
				Name:  "Zach Peters",
				Email: "zpeters@gmail.com",
			},
			Created:     post.Date,
			Description: truncate.Truncate(strippedHtml, 256, "...", truncate.PositionEnd),
		})
	}

	return feed.ToRss()
}

// parsePost splits a post into its title and date (from the front matter)
// and its markdown body. The front matter must give a title and a date
// between two --- lines.
func parsePost(src []byte) (string, time.Time, string, error) {
	src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(src, []byte("---\n")) {
		return "", time.Time{}, "", fmt.Errorf("no front matter")
	}
	meta, body, ok := bytes.Cut(src[len("---\n"):], []byte("\n---\n"))
	if !ok {
		return "", time.Time{}, "", fmt.Errorf("front matter is not closed with ---")
	}

	var m struct {
		Title string `yaml:"title"`
		Date  string `yaml:"date"`
	}
	if err := yaml.Unmarshal(meta, &m); err != nil {
		return "", time.Time{}, "", fmt.Errorf("front matter: %w", err)
	}
	if m.Title == "" {
		return "", time.Time{}, "", fmt.Errorf("front matter needs a title")
	}
	date, err := parseTimestamp(m.Date)
	if err != nil {
		return "", time.Time{}, "", fmt.Errorf("front matter: %w", err)
	}
	return m.Title, date, string(bytes.TrimSpace(body)), nil
}

func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}
