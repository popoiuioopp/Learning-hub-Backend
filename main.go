package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	function "github.com/popoiuioopp/Learning-hub-Backend/server/create"
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

	// Login or Join or Register
	// if login sucesss - createfc or host
	for ok := true; ok; ok = true{
		fmt.Println("Login or Register or Join Room : ")
		var usercmd string
		fmt.Scanln(&usercmd)
		if usercmd == "login"{
			//login
			if true{
				function.Create()
			}
		}else if usercmd == "register"{
			//regis
		}else{
			//join room
		}

	}

	var regisVerifyPass string
	fmt.Scanln(&regisVerifyPass)
}