package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var apiToken string

func main() {
	err := godotenv.Load("configs/api.env")
	if err != nil {
		fmt.Println(err)
		return
	}
	apiToken = os.Getenv("AUTH_TOKEN")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Mail",
		})
	})

	api := r.Group("/api", authMiddleware)
	api.GET("/mail", mail)
	api.GET("/mail/:id", mail)

	err = r.Run(":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func authMiddleware(c *gin.Context) {
	authHeader := c.Request.Header.Get("token")

	if authHeader == apiToken {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}

func mail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"mail": "mail",
	})
}
