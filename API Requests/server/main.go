package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

/* var books = []book{
	{ID: "1", Title: "Demon Slayer", Author: "Koyoharu Gotouge", Quantity: 10},
	{ID: "2", Title: "One Piece", Author: "Eiichiro Oda", Quantity: 5},
	{ID: "3", Title: "Bleach", Author: "Tite Kubo", Quantity: 7},
} */

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "library"
)

var Db *sql.DB
var err error

func init() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	Db, err = sql.Open("postgres", psqlconn)
	CheckError(err)
}

func main() {

	// populate the books data into the database
	/*for _, book := range books {
		_, err = db.Exec("INSERT INTO books (id, title, author, quantity) VALUES ($1, $2, $3, $4)", book.ID, book.Title, book.Author, book.Quantity)
		CheckError(err)
	}*/

	router := gin.Default()
	router.GET("/books", getBooksfromDb)
	router.POST("/books", createBookintoDb)
	router.GET("/books/:id", getBookbyIdfromDB)
	router.PATCH("/books/checkout/:id", checkoutBookfromDb)
	router.PATCH("/books/checkin/:id", checkInBookfromDb)
	router.Run("localhost:8080")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func getBooksfromDb(c *gin.Context) {

	//get data from database
	//return data
	rows, err := Db.Query("SELECT * FROM books")
	CheckError(err)
	defer rows.Close()

	// store data in books struct
	var books []book
	for rows.Next() {
		var id string
		var title string
		var author string
		var quantity int
		err = rows.Scan(&id, &title, &author, &quantity)
		CheckError(err)
		books = append(books, book{ID: id, Title: title, Author: author, Quantity: quantity})
	}

	// marshal to json
	c.IndentedJSON(http.StatusOK, books)
}

func createBookintoDb(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	_, err = Db.Exec("INSERT INTO books (id, title, author, quantity) VALUES ($1, $2, $3, $4)", newBook.ID, newBook.Title, newBook.Author, newBook.Quantity)
	CheckError(err)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookbyIdfromDB(c *gin.Context) {
	// get the id from the url
	id := c.Param("id")
	// get data from database
	// return data
	rows, err := Db.Query("SELECT * FROM books WHERE id = $1", id)
	CheckError(err)

	var tempBook book
	for rows.Next() {
		var id string
		var title string
		var author string
		var quantity int
		err = rows.Scan(&id, &title, &author, &quantity)
		CheckError(err)
		tempBook = book{ID: id, Title: title, Author: author, Quantity: quantity}
	}
	c.IndentedJSON(http.StatusOK, tempBook)
}

func checkoutBookfromDb(c *gin.Context) {
	id := c.Param("id")
	rows, err := Db.Query("SELECT * FROM books WHERE id = $1", id)
	CheckError(err)

	var tempBook book
	for rows.Next() {
		var id string
		var title string
		var author string
		var quantity int
		err = rows.Scan(&id, &title, &author, &quantity)
		CheckError(err)
		tempBook = book{ID: id, Title: title, Author: author, Quantity: quantity}
	}

	// check if the book is available or not
	if tempBook.Quantity == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available"})
		return
	}

	// update the quantity
	tempBook.Quantity = tempBook.Quantity - 1
	_, err = Db.Exec("UPDATE books SET quantity = $1 WHERE id = $2", tempBook.Quantity, tempBook.ID)
	CheckError(err)
	c.IndentedJSON(http.StatusOK, tempBook)
}

func checkInBookfromDb(c *gin.Context) {
	id := c.Param("id")
	rows, err := Db.Query("SELECT * FROM books WHERE id = $1", id)
	CheckError(err)
	var tempBook book
	for rows.Next() {
		var id string
		var title string
		var author string
		var quantity int
		err = rows.Scan(&id, &title, &author, &quantity)
		CheckError(err)
		tempBook = book{ID: id, Title: title, Author: author, Quantity: quantity}
	}

	// update the quantity
	tempBook.Quantity = tempBook.Quantity + 1
	_, err = Db.Exec("UPDATE books SET quantity = $1 WHERE id = $2", tempBook.Quantity, tempBook.ID)
	CheckError(err)
	c.IndentedJSON(http.StatusOK, tempBook)
}
