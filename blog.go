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
	"unicode"

	"github.com/aquilax/truncate"
	"github.com/goccy/go-yaml"
	"github.com/gorilla/feeds"
	"github.com/microcosm-cc/bluemonday"
)

type BlogPost struct {
	Date          time.Time
	FormattedDate string
	Title         string
	URL           string
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
		date, err := parseFilename(base)
		if err != nil {
			log.Printf("Skipping file %s due to parsing error: %v", base, err)
			continue
		}

		src, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Skipping file %s: %v", base, err)
			continue
		}
		title, content, err := parsePost(src)
		if err != nil {
			log.Printf("Skipping file %s: %v", base, err)
			continue
		}

		posts = append(posts, BlogPost{
			Date:          date,
			FormattedDate: FormatDate(date),
			Title:         title,
			URL:           slugify(title),
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

// parseFilename reads the post date from a file named <date>_<name>.md. The
// name after the first underscore is free text: the title comes from the
// front matter, so it can hold any characters.
func parseFilename(filename string) (time.Time, error) {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	dateTimeString, _, found := strings.Cut(base, "_")
	if !found {
		return time.Time{}, fmt.Errorf("filename must be <date>_<name>.md")
	}

	date, err := time.Parse("2006-01-02T15:04:05-0700", dateTimeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing datetime from filename: %v", err)
	}

	return date, nil
}

// parsePost splits a post into its title (from the front matter) and its
// markdown body. The front matter must be a title line between two --- lines.
func parsePost(src []byte) (string, string, error) {
	src = bytes.ReplaceAll(src, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(src, []byte("---\n")) {
		return "", "", fmt.Errorf("no front matter")
	}
	meta, body, ok := bytes.Cut(src[len("---\n"):], []byte("\n---\n"))
	if !ok {
		return "", "", fmt.Errorf("front matter is not closed with ---")
	}

	var m struct {
		Title string `yaml:"title"`
	}
	if err := yaml.Unmarshal(meta, &m); err != nil {
		return "", "", fmt.Errorf("front matter: %w", err)
	}
	if m.Title == "" {
		return "", "", fmt.Errorf("front matter needs a title")
	}
	return m.Title, string(bytes.TrimSpace(body)), nil
}

// slugify makes a URL path part from a title: lower case, letters and digits
// only, words joined by hyphens. Apostrophes are dropped, not hyphenated.
func slugify(title string) string {
	var b strings.Builder
	hyphen := false
	for _, r := range strings.ToLower(title) {
		switch {
		case r == '\'' || r == '’':
			// dropped
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			hyphen = false
		case !hyphen && b.Len() > 0:
			b.WriteByte('-')
			hyphen = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}
