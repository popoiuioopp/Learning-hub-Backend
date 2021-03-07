package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"

	"github.com/go-redis/redis"
)

// SQLHandler refers to the connection to the database.
type SQLHandler struct {
	Conn *sql.DB
}

var sqliteHandler SQLHandler

// User struct created when there is a signal to create user.
type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Check for the error
func checkErr(err error) int {
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func main() {
	redisClient := newClient()

	result, err := ping(redisClient)
	checkErr(err)
	fmt.Println(result)
}

func newClient() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})

	return redisClient
}

//Check connection between redis and the program.
func ping(client *redis.Client) (string, error) {
	result, err := client.Ping().Result()

	if checkErr(err) == 1 {
		return "", err
	}
	return result, nil
}
