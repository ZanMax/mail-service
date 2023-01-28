package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var apiToken string
var smtpHost string
var smtpPort string
var dbConnString string

func main() {
	err := godotenv.Load("configs/config.env")
	if err != nil {
		fmt.Println(err)
		return
	}
	apiToken = os.Getenv("AUTH_TOKEN")
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")
	dbConnString = os.Getenv("DB_CONNECTION_STRING")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Mail",
		})
	})

	api := r.Group("/api", authMiddleware)
	api.GET("/mail", mail)
	api.GET("/mail/:id", mail)
	api.POST("/mail/send/simple", mailSendSimple)

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
