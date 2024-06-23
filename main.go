package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


func main() {
	var router = gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, You created a Web App!"})
	})
	router.Static("/static", "./static")

	router.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.JSON(http.StatusOK, gin.H{"message": "Hello, " + name + "!"})
	})
	// Define a route to render an HTML template
	router.GET("/profile/:username", func(c *gin.Context) {
		username := c.Param("username")
		c.HTML(http.StatusOK, "profile.html", gin.H{"username": username})
	})
	router.Run(":8080")
}
