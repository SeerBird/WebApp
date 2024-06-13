package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"net/http"
)

var Socketio_Server *socketio.Server

func main() {
	var router = gin.Default()
	Socketio_Server = socketio.NewServer(nil)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, You created a Web App!"})
	})
	router.Static("/static", "./static")

	router.GET("/socket.io", socketHandler)
	router.POST("/socket.io", socketHandler)
	router.Handle("WS", "/socket.io", socketHandler)
	router.Handle("WSS", "/socket.io", socketHandler)

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

func socketHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, You created a Web App!"})
	Socketio_Server.OnConnect("connection", func(con socketio.Conn) error {
		fmt.Println("on connection")
		con.Join("chat")
		return nil
	})

	Socketio_Server.OnError("error", func(so socketio.Conn, err error) {
		fmt.Printf("[ WebSocket ] Error : %v", err.Error())
	})

	Socketio_Server.ServeHTTP(c.Writer, c.Request)
}
