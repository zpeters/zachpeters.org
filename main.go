// main is our entry point
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/hako/durafmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gin-gonic/gin"
)

type Post struct {
	FancyName     string
	Path          string
	Date          time.Time
	DateFormatted string
	ReadingTime   string
	Content       string
	WordCount     int
}

func main() {
	router := gin.Default()

	router.Use(gzip.Gzip(gzip.DefaultCompression))

	// Set cache control for static assets
	router.Use(func(c *gin.Context) {
		fmt.Println(c.Request.URL.Path)
		if strings.HasPrefix(c.Request.URL.Path, "/assets/") {
			c.Header("Cache-Control", "private, max-age=86400")
		}
		c.Next()
	})

	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.LoadHTMLGlob("templates/**/*")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
	})

	router.GET("/blog", func(c *gin.Context) {
		// TODO sort by file date
		var posts []Post

		var readingSpeed = 238

		files, err := os.ReadDir("./templates/blog/")
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		}

		caser := cases.Title(language.AmericanEnglish)

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			// get some file info
			info, err := file.Info()
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			}
			modTime := info.ModTime()

			// get some file content
			filePath := "./templates/blog/" + file.Name()
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			}

			// estimate reading time
			text := string(content)
			wordCount := len(strings.Fields(text))
			estimatedReadingMinutes := float64(wordCount) / float64(readingSpeed)
			estimatedTime := time.Duration(estimatedReadingMinutes) * time.Minute
			humanTime := durafmt.Parse(estimatedTime).String()

			// make our post object
			post := Post{
				FancyName:     caser.String(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))),
				Path:          file.Name(),
				Date:          modTime,
				DateFormatted: modTime.Format("Monday, January 02, 2006"),
				ReadingTime:   humanTime,
				WordCount:     wordCount,
			}
			posts = append(posts, post)
		}

		c.HTML(http.StatusOK, "bloglist.html", gin.H{
			"title": "Blog Entries",
			"posts": posts,
		})

	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Welcome",
		})
	})

	router.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.html", gin.H{
			"title": "Projects",
		})
	})

	router.GET("/recommendedapps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "recommendedapps.html", gin.H{
			"title": "Recommended Applications",
		})
	})

	log.Fatalln(router.Run(":8080"))
}
