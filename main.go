// main is our entry point
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"title": "Hello, World!",
		})
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

	router.GET("/changelog", func(c *gin.Context) {
		mdContent, err := os.ReadFile("./CHANGELOG.md")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		htmlContent := blackfriday.Run(mdContent)
		safeHtmlContent := template.HTML(htmlContent)

		c.HTML(http.StatusOK, "changelog.html", gin.H{
			"Content": safeHtmlContent,
		})
	})

	log.Fatalln(router.Run(":8080"))
}
