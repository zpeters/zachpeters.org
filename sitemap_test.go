package main

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSitemap(t *testing.T) {
	posts := []BlogPost{{URL: "recommended-apps", Date: time.Date(2026, 10, 18, 21, 30, 0, 0, time.UTC)}}
	out, err := GenerateSitemap(posts)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<loc>https://zachpeters.org/</loc>",
		"<loc>https://zachpeters.org/reads</loc>",
		"<loc>https://zachpeters.org/blog/recommended-apps</loc>",
		"<lastmod>2026-10-18</lastmod>",
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sitemap missing %q", want)
		}
	}
}
