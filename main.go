package main

import (
	"database/sql"
	"fmt"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/popoiuioopp/Learning-hub-Backend/cache"
)

// SQLHandler refers to the connection to the database.
type SQLHandler struct {
	Conn *sql.DB
}

var sqliteHandler SQLHandler
var forcreateuserid int
var redisHandler cache.RedisClient

// User struct created when there is a signal to create user.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int    `json:"id"`
}

//Check for the error

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

func CheckDeckExist(db *sql.DB, name string) int {
	var result int
	statement := `SELECT COUNT(*) FROM learninghub.Deck_instance where learninghub.Deck_instance.deckName = ?;`
	rows, err := db.Query(statement, name)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&result)
		fmt.Println(result)
		checkErr(err)
	}
	return result
}

//create flashard
func Createfc(db *sql.DB) {

	fmt.Println(">>>>>>>Create FlashCard")
	fmt.Printf("Create Deck????/(Y/N): ")
	var yesorno string
	fmt.Scanln(&yesorno)
	fmt.Println(yesorno)
	var deckname string
	if yesorno == "Y" {
		fmt.Println("Deckname:")
		fmt.Scanln(&deckname)

		if CheckDeckExist(db, deckname) == 0 {
			sqlStatement := `INSERT INTO Deck_instance(deckName, dateCreate) VALUES(?, NOW())`
			_, err := db.Exec(sqlStatement, deckname)

			checkErr(err)

		} else {
			fmt.Println("This Deck Name Already Used.")
			os.Exit(0)
		}

	} else {
		fmt.Println("See you next time")
		os.Exit(0)

	}

	var checkid int
	sqlStatement := `SELECT deckId FROM Deck_instance WHERE deckName = ? ORDER BY deckId DESC LIMIT 1 ` //check the lastest deckId and we will put it in the flashcard table
	rows, err := db.Query(sqlStatement, deckname)
	for rows.Next() {
		err = rows.Scan(&checkid)
		fmt.Println(checkid)
		checkErr(err)
	}
	checkErr(err)

	fmt.Println("Number of Flashcard : ") //let user choose
	var numfc int
	fmt.Scanln(&numfc)
	var slice []cache.FlashCard
	fmt.Println(slice)

	var temp cache.FlashCard
	for i := 0; i < numfc; i++ {
		fmt.Println("Term : ")
		fmt.Scanln(&temp.Term)
		fmt.Println("Definition : ")
		fmt.Scanln(&temp.Definition)
		slice = append(slice, temp)
	}

	var redisInstanceDeck cache.Deck
	for _, element := range slice {

		sqlStatement := `
		INSERT INTO Flashcard_instance(deckId,term,definition,userID)
		VALUES(?,?,?,?)
		`
		_, err := db.Exec(sqlStatement, checkid, element.Term, element.Definition, forcreateuserid)
		redisInstanceDeck.FlashCards = append(redisInstanceDeck.FlashCards, element)
		redisInstanceDeck.NoFC++
		checkErr(err)
	}

	sqlStatement = `select deck.deckName, deck.deckId 
	from Deck_instance as deck inner join Flashcard_instance as fc on
	deck.deckId = fc.deckId where deck.deckId = ? limit 1;`

	rows, err = db.Query(sqlStatement, checkid)
	for rows.Next() {
		err = rows.Scan(&redisInstanceDeck.DeckName, &redisInstanceDeck.DeckID)
		fmt.Println(checkid)
		checkErr(err)
	}
	cache.RedisAddDeck(redisHandler.Client, redisInstanceDeck)

}

// Create a User object and add to the database
func createUser(db *sql.DB) {
	fmt.Println("lets create your acc")
	var usercreate string
	var passcreate string

	fmt.Println("usernamecreate : ")
	fmt.Scanln(&usercreate)

	fmt.Println("passwordcreate : ")
	fmt.Scanln(&passcreate)

	sqlStatement := "insert into User(username, password) values(?, ?);" //
	_, err := db.Exec(sqlStatement, usercreate, passcreate)              // Execute the command
	checkErr(err)

}

func login(db *sql.DB) int {

	fmt.Println("lets login")

	fmt.Println("usernamelog : ")
	var username string
	fmt.Scanln(&username)

	fmt.Println("passwordlog : ")
	var password string
	fmt.Scanln(&password)
	sqlStatement := `SELECT username, password,userID FROM User WHERE username=? AND password=?`
	rows, err := db.Query(sqlStatement, username, password)
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
				fmt.Println("Successfully logged in")
			}

		}
	} else {
		fmt.Println("Cannot log in")
	}
	return queryResult[0].UserID
}

func ListDecks(db *sql.DB) {
	fmt.Println("==========")
	sqlStatement := "select deckid, deckName from Deck_instance;"
	rows, err := db.Query(sqlStatement)
	checkErr(err)
	for rows.Next() {
		var deckID string
		var deckName string
		err = rows.Scan(&deckID, &deckName)
		fmt.Printf("%s : %s\n", deckID, deckName)
	}
	fmt.Println("==========")
}

func main() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	checkErr(err)

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
	defer db.Close()

	redisHandler.Client = cache.NewClient()

	forcreateuserid = login(sqliteHandler.Conn)

	// Createfc(sqliteHandler.Conn)

	// result, err := cache.ReadDeck(redisHandler.Client, sqliteHandler.Conn, 1)
	// fmt.Println(result)

	result, err := cache.ReadDeck(redisHandler.Client, sqliteHandler.Conn, 1)
	checkErr(err)
	fmt.Println(result)

	var choice int
Loop:
	for {
		fmt.Printf("Please choose option:\n" +
			"1.)List all decks in the database \n" +
			"2.)Create Deck \n" +
			"3.)Check Deck Content \n" +
			"4.)Log out\n")
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			ListDecks(sqliteHandler.Conn)
		case 2:
			Createfc(sqliteHandler.Conn)
		case 3:
			fmt.Printf("Please Enter DeckId or 0 to exit:\n")
			fmt.Scanf("%d", &choice)
			if choice == 0 {
				continue
			} else {
				result, err := cache.ReadDeck(redisHandler.Client, sqliteHandler.Conn, 1)
				checkErr(err)
				fmt.Println(result)
			}
		default:
			fmt.Println("Logged out")
			break Loop
		}
	}
}
