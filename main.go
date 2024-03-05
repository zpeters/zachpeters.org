// main is our entry point
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/russross/blackfriday/v2"

	"github.com/gin-gonic/gin"
)

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

	router.LoadHTMLGlob("templates/**/*")

	blogPosts, err := LoadBlogPosts("blog/")
	if err != nil {
		panic(err)
	}

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Hello, World!",
		})
	})

	router.GET("/blog", func(c *gin.Context) {
		if err != nil {
			log.Printf("Could not load blog post list: %v", err)
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"message": "Could not load blog posts."})
			return
		}

		c.HTML(http.StatusOK, "blog.html", gin.H{
			"title": "Blog Posts",
			"posts": blogPosts,
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
			output := blackfriday.Run([]byte(foundPost.Content))
			foundPost.Content = string(output)

			c.HTML(http.StatusOK, "blog-post.html", gin.H{
				"title":   foundPost.Title,
				"date":    foundPost.Date,
				"content": template.HTML(foundPost.Content),
			})
		} else {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
		}
	})

	router.GET("/software", func(c *gin.Context) {
		c.HTML(http.StatusOK, "software-projects.html", gin.H{
			"title": "Software Projects",
		})
	})

	router.GET("/hardware", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hardware-projects.html", gin.H{
			"title": "Hardware Projects",
		})
	})

	router.GET("/recommendedapps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "recommendedapps.html", gin.H{
			"title": "Recommended Applications",
		})
	})

	router.GET("/articles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "articles.html", gin.H{
			"title": "Favorite Articles",
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
