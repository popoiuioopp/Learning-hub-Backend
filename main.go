package main

import (
	"database/sql"
	"fmt"

	"os"

	_ "github.com/go-sql-driver/mysql"

)

// SQLHandler refers to the connection to the database.
type SQLHandler struct {
	Conn *sql.DB
}

var sqliteHandler SQLHandler
var forcreateuserid int

// User struct created when there is a signal to create user.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int
}

//Check for the error

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())

	}
}

//create flashard
func Createfc() {

	type FlashCard struct {
		Term       string
		Definition string
	}

	fmt.Println(">>>>>>>Create FlashCard")
	fmt.Printf("Create Deck????/(Y/N): ")
	var yesorno string
	fmt.Scanln(&yesorno)
	fmt.Println(yesorno)
	if yesorno == "Y" {
		fmt.Println("Deckname:")
		var deckname string
		fmt.Scanln(&deckname)
		sqlStatement := `INSERT INTO Deck_instance(deckName) VALUES(?)`
		_, err := sqliteHandler.Conn.Exec(sqlStatement, deckname)

		checkErr(err)

	} else {
		fmt.Println("See you next time")
		os.Exit(0)

	}

	fmt.Println("Flashcard name : ")
	var namefc string
	var checkid int
	fmt.Scanln(&namefc)
	sqlStatement := `SELECT deckId FROM Deck_instance ORDER	BY deckId DESC LIMIT 1` //check the lastest deckId and we will put it in the flashcard table
	rows, err := sqliteHandler.Conn.Query(sqlStatement)
	for rows.Next() {
		err = rows.Scan(&checkid)
		fmt.Println(checkid)
		checkErr(err)
	}
	checkErr(err)

	fmt.Println("Number of Flashcard : ") //let user choose
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
		VALUES(?,?,?,?)
		`
		_, err := sqliteHandler.Conn.Exec(sqlStatement, checkid, element.Term, element.Definition, forcreateuserid)
		os.Exit(0)
		checkErr(err)
	}
}

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

func login() int {

	fmt.Println("usernamelog : ")
	var username string
	fmt.Scanln(&username)

	fmt.Println("passwordlog : ")
	var password string
	fmt.Scanln(&password)
	sqlStatement := `SELECT username, password,userID FROM User WHERE username=? AND password=?`
	rows, err := sqliteHandler.Conn.Query(sqlStatement, username, password)
	checkErr(err)
	var queryResult []User
	for rows.Next() { //check rows
		var tempUser User
		err = rows.Scan(&tempUser.Username, &tempUser.Password, &tempUser.UserID)
		checkErr(err)
		queryResult = append(queryResult, tempUser)
	}

	if len(queryResult) != 0 {

		fmt.Println(queryResult)
		for _, element := range queryResult { //if the username match the username in db then login success
			if element.Username == username && element.Password == password {
				fmt.Println("Successs loginnnnn")
			}

	} else {
		fmt.Println("Cannot log in")
	}
	return queryResult[0].UserID
}

func main() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	checkErr(err)

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
	defer db.Close()
	// createUser("DEARZA", "12345")
	// fmt.Println("Created successful")
	// createUser()
	forcreateuserid = login()
	Createfc()

	var quit string
	fmt.Scanln(&quit)
}

