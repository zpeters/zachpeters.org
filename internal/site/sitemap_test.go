package site

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateSitemap(t *testing.T) {
	posts := []BlogPost{{URL: "recommended-apps", Date: time.Date(2026, 10, 18, 21, 30, 0, 0, time.UTC)}}
	projects := []ProjectPage{{ID: "folder-cam"}}
	out, err := GenerateSitemap(posts, projects)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<loc>https://zachpeters.org/</loc>",
		"<loc>https://zachpeters.org/reads</loc>",
		"<loc>https://zachpeters.org/blog/recommended-apps</loc>",
		"<lastmod>2026-10-18</lastmod>",
		"<loc>https://zachpeters.org/projects/folder-cam</loc>",
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("sitemap missing %q", want)
		}
	}
}
