// main is our entry point
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

// version is the release version, set at build time with -ldflags
var version = "dev"

func main() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Set cache control for static assets
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "private, max-age=86400")
		}
		c.Next()
	})

	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.SetFuncMap(template.FuncMap{
		"contributions": StartContributions("zpeters", 6*time.Hour),
		"version":       func() string { return version },
	})
	router.LoadHTMLGlob("templates/**/*")

	blogPosts, err := LoadBlogPosts("blog/")
	if err != nil {
		panic(err)
	}

	reads, err := LoadReads("reads/")
	if err != nil {
		panic(err)
	}

	// Cheap liveness and readiness check; does not render any page
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Allow all crawlers, and point them at the sitemap
	router.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-agent: *\nAllow: /\n\nSitemap: https://zachpeters.org/sitemap.xml\n")
	})

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	router.GET("/", func(c *gin.Context) {
		// Get recent posts (limit to 5)
		recentPosts := blogPosts
		if len(recentPosts) > 5 {
			recentPosts = recentPosts[:5]
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"section":     "home",
			"recentPosts": recentPosts,
		})
	})

	router.GET("/blog", func(c *gin.Context) {
		c.HTML(http.StatusOK, "blog.html", gin.H{
			"section":  "blog",
			"title":    "Blog",
			"posts":    blogPosts,
			"timeline": NewPostTimeline(blogPosts, time.Now()),
		})
	})

	router.GET("/blog/:url", func(c *gin.Context) {
		postUrl := c.Param("url")

		var foundPost *BlogPost

		for _, post := range blogPosts {
			if post.URL == postUrl {
				foundPost = &post
				break
			}
		}

		if foundPost != nil {
			c.HTML(http.StatusOK, "blog-post.html", gin.H{
				"section":     "blog-post",
				"title":       foundPost.Title,
				"description": foundPost.Excerpt,
				"date":        foundPost.Date,
				"content":     foundPost.HTML,
			})
		} else {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
		}
	})

	router.GET("/software", func(c *gin.Context) {
		// Built from the project headings in the software-*.html templates,
		// so a new project shows up on the timeline without other edits
		timeline, err := LoadProjectTimeline("templates/pages/software-*.html", time.Now())
		if err != nil {
			log.Printf("Could not load project timeline: %v", err)
		}
		c.HTML(http.StatusOK, "software-projects.html", gin.H{
			"section":  "software",
			"title":    "Software",
			"timeline": timeline,
		})
	})

	router.GET("/hardware", func(c *gin.Context) {
		timeline, err := LoadProjectTimeline("templates/pages/hardware-*.html", time.Now())
		if err != nil {
			log.Printf("Could not load project timeline: %v", err)
		}
		c.HTML(http.StatusOK, "hardware-projects.html", gin.H{
			"section":  "hardware",
			"title":    "Hardware",
			"timeline": timeline,
		})
	})

	router.GET("/reads", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reads.html", gin.H{
			"section":  "reads",
			"title":    "Reads",
			"reads":    reads,
			"timeline": NewReadTimeline(reads, time.Now()),
		})
	})

	router.GET("/rss", func(c *gin.Context) {
		rss, err := GenerateRssFeed(blogPosts)
		if err != nil {
			log.Printf("Could not load rss feed: %v", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Could not load rss feed."})
			return
		}
		c.Writer.Header().Set("Content-Type", "application/rss+xml")
		c.String(http.StatusOK, rss)
	})

	log.Fatalln(router.Run(":8080"))
}
