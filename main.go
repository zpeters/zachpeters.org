package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
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

	router.LoadHTMLGlob("templates/*")

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "notfound.html", gin.H{})
	})

	router.GET("/markdown/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		markdown, err := os.ReadFile("markdown/" + filename + ".md")
		if err != nil {
			c.String(http.StatusNotFound, "file not found")
		}
		html := blackfriday.Run(markdown)
		c.Data(http.StatusOK, "text/html; charset=utf-8", html)
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
