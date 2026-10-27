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

// projectGlob matches the project pages, one file per project
const projectGlob = "projects/project-*.html"

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

	funcs := template.FuncMap{
		"contributions": StartContributions("zpeters", 6*time.Hour),
		"version":       func() string { return version },
	}
	// Shared includes and error pages live in templates/, the pages in pages/
	tmpl := template.Must(template.New("").Funcs(funcs).ParseGlob("templates/*/*"))
	tmpl = template.Must(tmpl.ParseGlob("pages/*.html"))
	router.SetHTMLTemplate(tmpl)

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

	router.GET("/sitemap.xml", func(c *gin.Context) {
		projects, err := LoadProjectPages(projectGlob)
		if err != nil {
			log.Printf("Could not load projects: %v", err)
		}
		sitemap, err := GenerateSitemap(blogPosts, projects)
		if err != nil {
			log.Printf("Could not build sitemap: %v", err)
			c.String(http.StatusInternalServerError, "could not build sitemap")
			return
		}
		c.Data(http.StatusOK, "application/xml; charset=utf-8", []byte(sitemap))
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

	router.GET("/projects", func(c *gin.Context) {
		// Built from the project headings in the project-*.html files,
		// so a new project shows up on the timeline without other edits
		timeline, err := LoadProjectTimeline(projectGlob, time.Now())
		if err != nil {
			log.Printf("Could not load project timeline: %v", err)
		}
		projects, err := LoadProjectPages(projectGlob)
		if err != nil {
			log.Printf("Could not load projects: %v", err)
		}
		c.HTML(http.StatusOK, "projects.html", gin.H{
			"section":  "projects",
			"title":    "Projects",
			"timeline": timeline,
			"projects": projects,
		})
	})

	router.GET("/projects/:id", func(c *gin.Context) {
		projects, err := LoadProjectPages(projectGlob)
		if err != nil {
			log.Printf("Could not load projects: %v", err)
		}
		for _, p := range projects {
			if p.ID == c.Param("id") {
				c.HTML(http.StatusOK, "project-page.html", gin.H{
					"section": "projects",
					"title":   p.Name,
					"project": p,
				})
				return
			}
		}
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	// The Software and Hardware pages are now Projects; keep their old links working
	router.GET("/software", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/projects") })
	router.GET("/hardware", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/projects") })
	router.GET("/software/:id", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/projects/"+c.Param("id"))
	})
	router.GET("/hardware/:id", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/projects/"+c.Param("id"))
	})

	router.GET("/colophon", func(c *gin.Context) {
		c.HTML(http.StatusOK, "colophon.html", gin.H{
			"title": "Colophon",
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
