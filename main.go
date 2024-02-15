// main is our entry point
package main

import (
	"fmt"
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
	TimeFormatted string
	ReadingTime   string
	Content       string
	WordCount     int
}

func filenameToFancyName(filename string) string {
	caser := cases.Title(language.AmericanEnglish)
	fancyName := caser.String(strings.TrimSuffix(filename, filepath.Ext(filename)))
	return fancyName
}

func getFileModTime(f os.DirEntry) (time.Time, error) {
	info, err := f.Info()
	if err != nil {
		return time.Time{}, err
	}
	modTime := info.ModTime()
	return modTime, nil
}

func getEstimatedReadingTime(filePath string) (string, int, error) {
	var readingSpeed = 238

	// get some file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", 0, err
	}

	// estimate reading time
	text := string(content)
	wordCount := len(strings.Fields(text))
	estimatedReadingMinutes := float64(wordCount) / float64(readingSpeed)
	estimatedTime := time.Duration(estimatedReadingMinutes) * time.Minute
	humanTime := durafmt.Parse(estimatedTime).String()

	return humanTime, wordCount, nil
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

		files, err := os.ReadDir("./templates/blog/")
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		}

		for _, file := range files {

			if file.IsDir() {
				continue
			}

			// get modtime
			modTime, err := getFileModTime(file)
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			}

			// get estimated reading time
			humanTime, wordCount, err := getEstimatedReadingTime("./templates/blog/" + file.Name())
			if err != nil {
				log.Println(err)
				c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			}

			// make our post object
			post := Post{
				FancyName:     filenameToFancyName(file.Name()),
				Path:          file.Name(),
				Date:          modTime,
				DateFormatted: modTime.Format("Monday, January 02, 2006"),
				TimeFormatted: modTime.Format("03:04 PM"),
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

	router.GET("/blog/:filename", func(c *gin.Context) {
		filename := c.Param("filename")

		// TODO make this work for file names not dir entries
		// modTime, err := getFileModTime(filename)
		// if err != nil {
		// 	log.Println(err)
		// 	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		// }
		// log.Println(modTime)

		// get estimated reading time
		humanTime, wordCount, err := getEstimatedReadingTime("./templates/blog/" + filename)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		}
		log.Println(humanTime)
		log.Println(wordCount)

		c.HTML(http.StatusOK, filename, gin.H{
			"title":    filenameToFancyName(filename),
			"metadata": fmt.Sprintf("%s", humanTime),
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
