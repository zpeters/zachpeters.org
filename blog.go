package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aquilax/truncate"
	"github.com/gorilla/feeds"
	"github.com/microcosm-cc/bluemonday"
)

type BlogPost struct {
	Date          time.Time
	FormattedDate string
	Title         string
	URL           string
	Content       string
	Excerpt       string // first sentence, plain text
}

func LoadBlogPosts(postsDir string) ([]BlogPost, error) {
	var posts []BlogPost

	files, err := filepath.Glob(filepath.Join(postsDir, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		base := filepath.Base(file)
		date, title, err := parseFilename(base)
		if err != nil {
			log.Printf("Skipping file %s due to parsing error: %v", base, err)
			continue
		}

		contentBytes, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Skipping file %s: %v", base, err)
			continue
		}
		content := string(contentBytes)

		posts = append(posts, BlogPost{
			Date:          date,
			FormattedDate: FormatDate(date),
			Title:         title,
			URL:           strings.ToLower(strings.ReplaceAll(title, " ", "-")),
			Content:       content,
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
		strippedHtml := p.Sanitize(post.Content)
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

func parseFilename(filename string) (time.Time, string, error) {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	underscoreIndex := strings.LastIndex(base, "_")
	if underscoreIndex == -1 {
		return time.Time{}, "", fmt.Errorf("filename does not contain a title part")
	}

	dateTimeString := string(base[:underscoreIndex])
	title := strings.ReplaceAll(base[underscoreIndex+1:], "-", " ")

	date, err := time.Parse("2006-01-02T15:04:05-0700", dateTimeString)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("error parsing datetime from filename: %v", err)
	}

	return date, title, nil
}

func FormatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}
