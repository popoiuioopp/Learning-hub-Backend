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

// func Createfc() {

// 	type FlashCard struct {
// 		Term       string
// 		Definition string
// 	}

// 	fmt.Println(">>>>>>>Create FlashCard")

// 	fmt.Println("Flashcard name : ")
// 	var namefc string
// 	fmt.Scanln(&namefc)

// 	fmt.Println("Number of Flashcard : ")
// 	var numfc int
// 	fmt.Scanln(&numfc)
// 	var slice []FlashCard
// 	fmt.Println(slice)

// 	var temp FlashCard
// 	for i := 0; i < numfc; i++ {
// 		fmt.Println("Term : ")
// 		fmt.Scanln(&temp.Term)
// 		fmt.Println("Definition : ")
// 		fmt.Scanln(&temp.Definition)
// 		slice = append(slice, temp)
// 	}
// 	fmt.Println(slice)
// 	fmt.Println(len(slice))

// 	for _, element := range slice {
// 		// fmt.Println(index, element.Term)
// 		sqlStatement := `
// 		INSERT INTO Flashcard_instance(deckId,term,definition,userID)
// 		VALUES(1,?,?,1)
// 		`
// 		_, err := sqliteHandler.Conn.Exec(sqlStatement, element.Term, element.Definition)
// 		checkErr(err)
// 	}
// }

// Create a User object and add to the database
func createUser() {
	var usercreate string
	var passcreate string

	fmt.Println("usernamecreate : ")
	fmt.Scanln(&usercreate)

	fmt.Println("passwordcreate : ")
	fmt.Scanln(&passcreate)

	sqlStatement := "insert into User(username, password) values(?, ?);"    //
	_, err := sqliteHandler.Conn.Exec(sqlStatement, usercreate, passcreate) // Execute the command
	checkErr(err)
	fmt.Println("Insert dai na")

}

func login() {
	fmt.Println("usernamelog : ")
	var username string
	fmt.Scanln(&username)

	fmt.Println("passwordlog : ")
	var password string
	fmt.Scanln(&password)
	sqlStatement := `SELECT username, password FROM User WHERE username=? AND password=?`
	rows, err := sqliteHandler.Conn.Query(sqlStatement, username, password)
	checkErr(err)
	var queryResult []User
	for rows.Next() {
		var tempUser User
		err = rows.Scan(&tempUser.Username, &tempUser.Password)
		checkErr(err)
		queryResult = append(queryResult, tempUser)
	}

	if len(queryResult) != 0 {

		fmt.Println(queryResult)
		for _, element := range queryResult {
			if element.Username == username && element.Password == password {
				fmt.Println("Successs loginnnnn IMPORT BOOSSSSS")
			} else {
				fmt.Println("Noooooo")
			}
		}
	} else {
		fmt.Println("Cannot log in")
	}

}

func main() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	checkErr(err)

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
	defer db.Close()
	// createUser("DEARZA", "12345")
	// fmt.Println("Created successful")
	Createfc()
	// createUser()

	var quit string
	fmt.Scanln(&quit)
}
