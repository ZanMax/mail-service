package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
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
	reqErr := c.BindJSON(&simpleMail)
	if reqErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad request",
		})
		return
	}

	db, dbErr := sql.Open("mysql", dbConnString)
	defer db.Close()

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	res, errSQL := db.Query("SELECT password from emails where email=?", simpleMail.From)
	if errSQL != nil {
		log.Fatal(errSQL)
	}
	defer res.Close()
	var pwd string
	if res.Next() {
		res.Scan(&pwd)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not found",
		})
		return
	}
	templateFile := "test.html"
	params := struct {
		Title   string
		Message string
	}{
		Title:   "Title",
		Message: "bodyMsg",
	}
	var files []string
	go MailSender(simpleMail.From, simpleMail.To, pwd, simpleMail.Subject, simpleMail.Message, files, templateFile, params)

	c.JSON(http.StatusOK, gin.H{
		"status": "sent",
	})
}

func mail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func mailEvent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
