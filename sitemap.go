package main

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

// GenerateSitemap lists the fixed pages and every blog post. A post's
// last-modified date is its post date.
func GenerateSitemap(posts []BlogPost) (string, error) {
	set := urlSet{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9"}
	for _, path := range []string{"/", "/blog", "/software", "/hardware", "/reads"} {
		set.URLs = append(set.URLs, sitemapURL{Loc: siteURL + path})
	}
	for _, p := range posts {
		set.URLs = append(set.URLs, sitemapURL{
			Loc:     siteURL + "/blog/" + p.URL,
			LastMod: p.Date.Format(time.DateOnly),
		})
	}
	out, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(out) + "\n", nil
}
