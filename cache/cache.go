package cache

import (
	"database/sql"
	"encoding/json"
	"fmt"

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
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   int    `json:"id"`
}

type Deck struct {
	DeckID     int         `json:"id"`
	DeckName   string      `json:"name"`
	NoFC       int         `json:"number"`
	FlashCards []FlashCard `json:"FlashCards"`
}

//FlashCard struct for containing the FlashCard
type FlashCard struct {
	Term       string `json:"term"`
	Definition string `json:"defintion"`
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

func RedisAddDeck(client *redis.Client, deck Deck) {
	deckKey := fmt.Sprintf("deck,%d", deck.DeckID)

	temp, err := json.Marshal(deck)
	checkErr(err)

	_, err = client.Get(deckKey).Result()

	if err == redis.Nil {
		err = client.Set(deckKey, temp, 0).Err()
		checkErr(err)
	}

}

func ReadDeck(client *redis.Client, db *sql.DB, id int) (string, error) {
	result, err := client.Get(fmt.Sprintf("deck,%d", id)).Result()
	var redisInstanceDeck Deck
	if err == redis.Nil {
		sqlStatement := `select deck.deckid, deck.deckName, fc.term, fc.definition from Deck_instance as deck inner join Flashcard_instance as fc 
		on deck.deckId = fc.deckId inner join User as user on fc.userID = user.userID
		 where deck.deckId = ?;`

		rows, err := db.Query(sqlStatement, id)
		checkErr(err)

		for rows.Next() {
			var tempFlashCard FlashCard
			err = rows.Scan(&redisInstanceDeck.DeckID, &redisInstanceDeck.DeckName, &tempFlashCard.Term, &tempFlashCard.Definition)
			checkErr(err)
			redisInstanceDeck.FlashCards = append(redisInstanceDeck.FlashCards, tempFlashCard)
			redisInstanceDeck.NoFC++
		}
		temp, err := json.Marshal(redisInstanceDeck)
		result = string(temp)
	}

	return result, nil
}
