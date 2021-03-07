package create

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
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
}

// create Flashcard and store in database
func Create() {

	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
	CheckErr(err)
	sqliteHandler.Conn = db

	Createfc()

	var regisVerifyPass string
	fmt.Scanln(&regisVerifyPass)
}

func Createfc() {

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
		CheckErr(err)
	}
}
