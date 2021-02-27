package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// SQLHandler refers to the connection to the database.
type SQLHandler struct {
	Conn *sql.DB
}

var sqliteHandler SQLHandler

// check error
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

// create Flashcard and store in database
func main() {
	// limiter := rate.NewLimiter(10, 1)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	if !limiter.Allow() {
	// 		w.WriteHeader(http.StatusTooManyRequests)
	// 		return
	// 	}
	// 	fmt.Fprintln(w, "Hello, World")
	// })
	// log.Fatal(http.ListenAndServe(":8080", nil))
	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	checkErr(err)
	sqliteHandler.Conn = db

	create()

	var regisVerifyPass string
	fmt.Scanln(&regisVerifyPass)
}

func create() {

	type FlashCard struct {
		Term       string
		Definition string
	}

	fmt.Println(">>>>>>>Create FlashCard")

	fmt.Println("Flashcard name : ")
	var namefc string
	fmt.Scanln(&namefc)

	fmt.Println("Number of Flashcard : ")
	var numfc int
	fmt.Scanln(&numfc)
	var slice []FlashCard
	fmt.Println(slice)

	var temp FlashCard
	for i := 0; i < numfc; i++ {
		fmt.Println("Term : ")
		fmt.Scanln(&temp.Term)
		fmt.Println("Definition : ")
		fmt.Scanln(&temp.Definition)
		slice = append(slice, temp)
	}
	fmt.Println(slice)
	fmt.Println(len(slice))

	for _, element := range slice {
		// fmt.Println(index, element.Term)
		sqlStatement := `
		INSERT INTO Flashcard_instance(deckId,term,definition,userID)
		VALUES(1,?,?,1)
		`
		_, err := sqliteHandler.Conn.Exec(sqlStatement, element.Term, element.Definition)
		checkErr(err)
	}
}
