package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type SimpleMail struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
}

func getDomainFromEmail(email string) string {
	domain := strings.Split(email, "@")
	return domain[1]
}

func getMailCredentials(email string) string {
	domain := getDomainFromEmail(email)
	fmt.Println(domain)
	return ""
}

func mailSendSimple(c *gin.Context) {
	var simpleMail SimpleMail
	err := c.BindJSON(&simpleMail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}
	fmt.Println(simpleMail)
	c.JSON(http.StatusOK, gin.H{
		"status": "sent",
	})
}
