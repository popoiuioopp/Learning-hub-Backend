package main

import (
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"database/sql"

	"github.com/go-redis/redis"
)

// SQLHandler refers to the connection to the database.
type SQLHandler struct {
	Conn *sql.DB
}

// RedisClient is a variable that contains redis client
type RedisClient struct {
	Client *redis.Client
}

var redisHandler RedisClient

var sqliteHandler SQLHandler

// User struct created when there is a signal to create user.
type User struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//FlashCard struct for containing the FlashCard
type FlashCard struct {
	FlashCardID int    `json:"id"`
	Term        string `json:"term"`
	Definition  string `json:"definition"`
}

func main() {
	redisHandler.Client = NewClient()

	result, err := ping(redisHandler.Client)
	checkErr(err)
	fmt.Println(result)

	err = addUsr(redisHandler.Client, 1, "Boooosss", "password")
	checkErr(err)
	result, err = checkUsr(redisHandler.Client, 1)
	checkErr(err)

	fmt.Println(result)
}

//Check for the error
func checkErr(err error) int {
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

//NewClient will create new Redis client
func NewClient() *redis.Client {
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

func addUsr(client *redis.Client, userID int, userName string, userPwd string) error {
	userKey := fmt.Sprintf("usr,%d", userID)
	temp, err := json.Marshal(User{UserID: userID, Username: userName, Password: userPwd})
	_, err = client.Get(userKey).Result()
	if err == redis.Nil {
		err = client.Set(userKey, temp, 0).Err()

		checkErr(err)
		return err
	}
	fmt.Println("User already exits")
	return err
}

func checkUsr(client *redis.Client, userID int) (string, error) {

	val, err := client.Get(fmt.Sprintf("usr,%d", userID)).Result()
	checkErr(err)

	return val, err
}
