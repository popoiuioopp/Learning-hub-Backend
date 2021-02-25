// package main

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/popoiuioopp/Learning-hub-Backend/server"
// )

// func EqualFold(s, t string) bool {
// 	return strings.EqualFold(s, t)
// }

// func login() {
// 	fmt.Println("Enter Code or login or register: ")

// 	// var then variable name then variable type
// 	var first string

// 	// Taking input from user
// 	fmt.Scanln(&first)
// 	if EqualFold(first, "login") {
// 		fmt.Println("Enter username : ")
// 		var username string
// 		fmt.Scanln(&username)

// 		fmt.Println("Enter password : ")
// 		var password string
// 		fmt.Scanln(&password)

// 		// if login sucess

// 	} else if EqualFold(first, "register") {
// 		fmt.Println(">>>>>>>Register")
// 		fmt.Println("Username : ")
// 		var regisUsername string
// 		fmt.Scanln(&regisUsername)

// 		fmt.Println(">>>>>>>Register")
// 		fmt.Println("Password : ")
// 		var regisPass string
// 		fmt.Scanln(&regisPass)

// 		// fmt.Println(">>>>>>>Register")
// 		// fmt.Println("Verify Password :")
// 		// var regisVerifyPass string
// 		// fmt.Scanln(&regisVerifyPass)

// 		// for regisPass != regisPass {
// 		// 	fmt.Println(">>>>>>>Register : Password do now match")
// 		// 	fmt.Println("Verify Password :")
// 		// 	var regisVerifyPass string
// 		// 	fmt.Scanln(&regisVerifyPass)
// 		// }
// 	} else {
// 		fmt.Println("http://localhost/" + first)
// 	}
// }
// func main() {
// 	/*
// 		// fmt.Println("Hello World! Tangy Ma Laew")
// 		// Println function is used to
// 		// display output in the next line
// 		for {
// 			login()
// 			if false {
// 				break
// 			}
// 		}
// 		var wait string
// 		fmt.Scanln(&wait)
// 	*/
// 	// server.Server()
// }
package main

import (
	"database/sql"
	"fmt"
<<<<<<< Updated upstream
	"strings"
)

var check bool

func EqualFold(s, t string) bool {
	return strings.EqualFold(s, t)
}

func login() {
	fmt.Println("Enter Code or login or register: ")

	// var then variable name then variable type
	var first string

	// Taking input from user
	fmt.Scanln(&first)
	if EqualFold(first, "login") {
		fmt.Println("Enter username : ")
		var username string
		fmt.Scanln(&username)

		fmt.Println("Enter password : ")
		var password string
		fmt.Scanln(&password)

		// if login sucess
		fmt.Println("Login Sucess!!!!")

		fmt.Println("Create Flashcard or choose Flashcard : ")
		var userInput string
		fmt.Scanln(&userInput)

	} else if EqualFold(first, "register") {
		fmt.Println(">>>>>>>Register")
		fmt.Println("Username : ")
		var regisUsername string
		fmt.Scanln(&regisUsername)

		fmt.Println(">>>>>>>Register")
		fmt.Println("Password : ")
		var regisPass string
		fmt.Scanln(&regisPass)

		// fmt.Println(">>>>>>>Register")
		// fmt.Println("Verify Password :")
		// var regisVerifyPass string
		// fmt.Scanln(&regisVerifyPass)

		// for regisPass != regisPass {
		// 	fmt.Println(">>>>>>>Register : Password do now match")
		// 	fmt.Println("Verify Password :")
		// 	var regisVerifyPass string
		// 	fmt.Scanln(&regisVerifyPass)
		// }

		fmt.Println("Register Sucess!!!!")
	} else {
		fmt.Println("Connecting to Room : " + first)
		check = true
=======

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
>>>>>>> Stashed changes
	}
}

// create Flashcard and store in database
func main() {
<<<<<<< Updated upstream
	// fmt.Println("Hello World! Tangy Ma Laew")
	// Println function is used to
	// display output in the next line
	for {
		login()
		if check {
			break
		}
	}
	var wait string
	fmt.Scanln(&wait)
=======
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
>>>>>>> Stashed changes
}
