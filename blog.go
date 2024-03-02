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
	Date    time.Time
	Title   string
	URL     string
	Content string // You might load this only when needed
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
			Date:    date,
			Title:   title,
			URL:     strings.ToLower(strings.ReplaceAll(title, " ", "-")),
			Content: content,
		})

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})
	}

	return posts, nil
}

func parseFilename(filename string) (time.Time, string, error) {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	parts := strings.SplitN(base, "-", 4)
	if len(parts) < 4 {
		return time.Time{}, "", fmt.Errorf("filename does not match expected format: YYYY-MM-DD-title")
	}

	dateString := strings.Join(parts[:3], "-")
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return time.Time{}, "", err
	}

	title := strings.Join(parts[3:], "-")
	title = strings.Replace(title, "-", " ", -1)

	return date, title, nil
}
