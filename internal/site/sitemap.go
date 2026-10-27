package site

import (
	"encoding/xml"
	"time"
)

const siteURL = "https://zachpeters.org"

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type urlSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// GenerateSitemap lists the fixed pages, every blog post and every project
// page. A post's last-modified date is its post date; projects have no
// date, so they have none.
func GenerateSitemap(posts []BlogPost, projects []ProjectPage) (string, error) {
	set := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, path := range []string{"/", "/blog", "/projects", "/reads", "/colophon"} {
		set.URLs = append(set.URLs, sitemapURL{Loc: siteURL + path})
	}
	for _, p := range posts {
		set.URLs = append(set.URLs, sitemapURL{
			Loc:     siteURL + "/blog/" + p.URL,
			LastMod: p.Date.Format(time.DateOnly),
		})
	}
	for _, p := range projects {
		set.URLs = append(set.URLs, sitemapURL{Loc: siteURL + "/projects/" + p.ID})
	}
	out, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(out) + "\n", nil
}
