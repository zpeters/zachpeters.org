package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type BlogPost struct {
	Date          time.Time
	FormattedDate string
	Title         string
	URL           string
	Content       string // You might load this only when needed
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
			continue
		}
		content := string(contentBytes)

		posts = append(posts, BlogPost{
			Date:          date,
			FormattedDate: FormatDate(date),
			Title:         title,
			URL:           strings.ToLower(strings.ReplaceAll(title, " ", "-")),
			Content:       content,
		})

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})
	}

	return posts, nil
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
	baseFormat := t.Format("Monday, Jan 2, 2006")
	timeFormat := t.Format("at 03:04PM")
	return fmt.Sprintf("%s %s", baseFormat, timeFormat)
}
