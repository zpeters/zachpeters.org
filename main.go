package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday/v2"
)

func main() {
	router := gin.Default()

	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")

	router.LoadHTMLGlob("templates/*")

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

	router.GET("/playground", func(c *gin.Context) {
		c.HTML(http.StatusOK, "playground.html", gin.H{
			"title": "CSS Playground",
		})
	})

	router.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.html", gin.H{
			"title": "Projects",
		})
	})

	log.Fatalln(router.Run(":8080"))
}
