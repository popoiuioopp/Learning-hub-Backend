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

// User struct created when there is a signal to create user.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Check for the error
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func main() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	checkErr(err)

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
	defer db.Close()

}
