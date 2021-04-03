package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	num_user int
}

type SQLHandler struct {
	Conn *sql.DB
}

type RedisClient struct {
	Client *redis.Client
}

type Flashcard struct {
	FcID       int
	Term       string
	Definition string
	DeckID     int
}

type Deck struct {
	DeckID   int
	DeckName string
	FcArray  *[]Flashcard
}

var sqliteHandler SQLHandler
var redisHandler RedisClient

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
		num_user: 0,
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args[1])
		case CMD_JOIN:
			s.join(cmd.client, cmd.args[1])
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_CROOM:
			s.croom(cmd.client, cmd.args[1])
		// case CMD_CREATEFC:
		// 	s.createfc(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_CUSER:
			s.cuser(cmd.client)
		case CMD_SRD:
			s.setroomdeck(cmd.client, cmd.args[1])
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())
	s.num_user += 1

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		status:   "0",
		commands: s.commands,
		score:    0,
		vaild:    false,
		no_ques:  0,
	}

	c.readInput()
}

func NewDBConn() {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("mysql", "learninghub:FgTQTzNM62cC63K@tcp(143.198.204.127:3306)/learninghub")
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("Connected to database")
	sqliteHandler.Conn = db
}

//NewClient will create new Redis client
func NewRedisConn() {
	fmt.Println("Connecting to Redis....")
	redisClient := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // default DB
	})
	redisHandler.Client = redisClient
	_, err := redisHandler.Client.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to Redis")
		redisHandler.Client.FlushAll().Result()
	}
}

func (s *server) nick(c *client, nick string) {
	c.msg(fmt.Sprint(nick))
	if len(nick) == 0 {
		c.msg(fmt.Sprintf("Invalid syntax"))
	}

	c.nick = nick
	c.msg(fmt.Sprintf("All right, We will call you %s", nick))
}

func (s *server) croom(c *client, roomName string) {

	//c.msg(fmt.Sprint(roomName))
	// if len(roomName) == 0 {
	// 	c.msg(fmt.Sprintf("Invalid syntax"))
	// }

	c.msg(fmt.Sprintf("croom %s", roomName))
	r, ok := s.rooms[roomName]

	if !ok {
		ip := c.conn.RemoteAddr().String()
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
			host:    ip, // ip of client
			status:  false,
			answer:  "",
			no_fc:   0,
		}

		s.rooms[roomName] = r

	} else {
		fmt.Println("Room already existed!...")
		return
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	c.msg(fmt.Sprintf("welcome to %s\nIf you want to play the game type /ready", roomName))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {
		c.msg(fmt.Sprintf("Room is not available!"))
		return
	}
	if r.status == false {
		r.members[c.conn.RemoteAddr()] = c

		s.quitCurrentRoom(c)
		c.room = r

		r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

		c.msg(fmt.Sprintf("welcome to %s\nIf you want to play the game type /ready", roomName))

	} else {
		c.msg(fmt.Sprintf("game already start"))
	}

}

//////list all deck//////////////////////
func ListDecks(c *client) {
	sqlStatement := "select deckid, deckName from Deck_instance;"
	rows, err := sqliteHandler.Conn.Query(sqlStatement)
	if err != nil {
		return
	}
	for rows.Next() {
		var deckID string
		var deckName string
		err = rows.Scan(&deckID, &deckName)
		if err != nil {
			return
		}
		c.msg(fmt.Sprintf("%s : %s\n", deckID, deckName))
	}
}

////////////Check that deck name is already existed or not////////////////////////////
func checkDeckExist(deckname string) int {
	var result int
	statement := `SELECT COUNT(*) FROM learninghub.Deck_instance 
	where learninghub.Deck_instance.deckName = ?;`
	rows, err := sqliteHandler.Conn.Query(statement, deckname)
	if err != nil {
		return 1
	}
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return 1
		}
	}
	return result
}

////////////Check that deck name is already existed or not////////////////////////////
func checkDeckIdExist(deckid int) int {
	var result int
	statement := `SELECT COUNT(*) FROM learninghub.Deck_instance 
	where learninghub.Deck_instance.deckId = ?;`
	rows, err := sqliteHandler.Conn.Query(statement, deckid)
	if err != nil {
		return 0
	}
	for rows.Next() {
		err = rows.Scan(&result)
		if err != nil {
			return 0
		}
	}
	return 1
}

////////////time measure//////////////////////////
func timesup(c *client, msg string) {
	c.msg(msg)
	time.Sleep(3 * time.Second)
}

//check the lastest deckId and we will put it in the flashcard table
func checkDeckId(deckname string) (int, error) {
	var checkid int
	sqlStatement := `SELECT deckId FROM Deck_instance 
		WHERE deckName = ? ORDER BY deckId DESC LIMIT 1 `
	rows, err := sqliteHandler.Conn.Query(sqlStatement, deckname)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&checkid)
		if err != nil {
			return 0, err
		}
	}
	return checkid, nil
}

func createDeck(c *client, deckname string) bool {
	if checkDeckExist(deckname) == 0 {
		sqlStatement := `INSERT INTO Deck_instance(deckName, dateCreate) VALUES(?, NOW())`
		_, err := sqliteHandler.Conn.Exec(sqlStatement, deckname)

		if err != nil {
			return false
		}
		c.msg(fmt.Sprintf("be able to create"))
		return true
	} else {
		c.msg(fmt.Sprintf("This Deck Name Already Used. You will be return to lobby"))
		return false
	}
}

func createfc(c *client, listFC []Flashcard) {
	sqlStatement := `
		INSERT INTO Flashcard_instance(deckId,term,definition)
		VALUES(?,?,?)`
	for _, item := range listFC {
		// c.msg(fmt.Sprintf("%d, %s, %s\n", item.deckID, item.term, item.definition))
		_, err := sqliteHandler.Conn.Exec(sqlStatement, item.DeckID, item.Term, item.Definition)
		if err != nil {
			return
		}
	}
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("There are currently %d rooms\navailable rooms: %s", len(rooms),
		strings.Join(rooms, ", ")))

}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) cuser(c *client) {
	// see all number of user in the server or room
	return
}

func (s *server) setroomdeck(c *client, deckid string) {
	// set room.deckid

	if c.conn.RemoteAddr().String() == c.room.host {
		deck_id, err := strconv.Atoi(deckid)
		if err != nil {
			c.msg(fmt.Sprintf("%s\n", err))
		}
		unmar, err := redisHandler.Client.Get(deckid).Result()
		if err == redis.Nil {
			// deckStructRedis, err := redisHandler.Client.Get(deckid).Result()
			sqlStatement := `select Deck_instance.deckName, Deck_instance.deckId, 
		Flashcard_instance.flashcardID, Flashcard_instance.Term, 
		Flashcard_instance.definition from Flashcard_instance 
		inner join Deck_instance
		on Flashcard_instance.deckId = Deck_instance.deckId
		where Deck_instance.deckId = ?; `
			rows, err := sqliteHandler.Conn.Query(sqlStatement, deck_id)
			if err != nil {
				return
			}
			var fcArray []Flashcard
			for rows.Next() {
				var tempFC Flashcard
				tempFC.DeckID = deck_id
				err := rows.Scan(&c.room.deck.DeckName, &c.room.deck.DeckID,
					&tempFC.FcID, &tempFC.Term, &tempFC.Definition)
				fcArray = append(fcArray, tempFC)
				if err != nil {
					return
				}
				c.room.deck.FcArray = &fcArray
			}

			var jsonData []byte
			jsonData, err = json.Marshal(c.room.deck)
			if err != nil {
				return
			}
			redisHandler.Client.Set(fmt.Sprintf("%d", c.room.deck.DeckID), string(jsonData), 0)
			jsonData = nil
		} else {
			b := []byte(unmar)
			deck := &Deck{}
			err = json.Unmarshal(b, deck)
			if err != nil {
				return
			}
			c.room.deck = *deck
		}
		c.room.no_fc = len(*c.room.deck.FcArray)
		c.msg(fmt.Sprintf("This room have these detail\nDeckID=%d\nDeckName=%s\n",
			c.room.deck.DeckID, c.room.deck.DeckName))
		for _, item := range *c.room.deck.FcArray {
			c.msg(fmt.Sprintf("%d, %d, %s, %s\n", item.FcID, item.DeckID, item.Term,
				item.Definition))
		}
		return
	} else {
		c.msg(fmt.Sprintf("You don't have the permission to change the room deck!"))
		return
	}
}
