package authen

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

func Login() {
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

// func main() {

// 	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(139.59.106.148:3306)/learninghub")
// 	checkErr(err)

// 	fmt.Println("Connected to database")
// 	sqliteHandler.Conn = db
// 	defer db.Close()
// 	// createUser("DEARZA", "12345")
// 	// fmt.Println("Created successful")
// 	//login()
// 	// createUser()

// 	var quit string
// 	fmt.Scanln(&quit)
// }
