package main

import "github.com/gin-gonic/gin"

// Ping get /api/ping
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
