package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// send a request to the server to get all the books

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func main() {

	router := gin.Default()
	router.GET("/getbooksclient", GetBooks)
	router.POST("/createbookclient", insertBook)
	router.Run("localhost:8081")

}

func GetBooks(c *gin.Context) {

	baseURL := "http://localhost:8080/books"
	resp, err := http.Get(baseURL)
	if err != nil {
		panic("Server is not running")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(string(respBody))

	// unmarshal the response
	var books []book
	err = json.Unmarshal(respBody, &books)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// send the response
	c.IndentedJSON(http.StatusOK, books)
}

func insertBook(c *gin.Context) {
	// marshal the book
	baseURL := "http://localhost:8080/books"

	data := c.Request.Body

	// post the gin context to the server
	resp, err := http.Post(baseURL, "application/json", data)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	// fmt.Println(string(respBody))

	// unmarshal the response
	var book book
	err = json.Unmarshal(respBody, &book)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// send the response
	c.IndentedJSON(http.StatusOK, book)

}
