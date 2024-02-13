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
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home Page",
		})
	})

	log.Fatalln(router.Run(":8080"))
}
