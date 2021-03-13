package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func readMsg(c *client) {
	c.msg(fmt.Sprintf("hello in readMsg"))

	for {
		c.msg(fmt.Sprintf("loop"))
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.msg(fmt.Sprintf("loop1"))
		if err != nil {
			return
		}
		c.msg(fmt.Sprintf(msg))
		c.msg(fmt.Sprintf("loop2"))
	}

}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
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
		case CMD_START:
			s.startGame(cmd.client)
		case CMD_CREATEFC:
			s.createfc(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_LOGIN:
			s.login(cmd.client, cmd.args[1], cmd.args[2])
		case CMD_REGIS:
			s.regis(cmd.client, cmd.args[1], cmd.args[2], cmd.args[3])
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		status:   "0",
		commands: s.commands, 
	}

	c.readInput()
}

func (s *server) nick(c *client, nick string) {
	c.nick = nick
	c.msg(fmt.Sprintf("all right, I will call you %s", nick))
}

func (s *server) croom(c *client, roomName string) {

	c.msg(fmt.Sprintf("croom %s", roomName))
	r, ok := s.rooms[roomName]

	if !ok {

		c.msg(fmt.Sprintf("Select Your Flashcard..."))
		// Put Select Flashcard Function From Boss Here

		r = &room{
			name:      roomName,
			members:   make(map[net.Addr]*client),
			flashcard: "Flashcard", // Selected Flashcard
			host:      "host_id",   // Host ID
		}
		s.rooms[roomName] = r
	} else {
		fmt.Println("Room already existed!...")
		return
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	// r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) join(c *client, roomName string) {
	r, ok := s.rooms[roomName]
	if !ok {

		c.msg(fmt.Sprintf("Room is not available!"))
		return
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) startGame(c *client) {
	//start game
}

func (s *server) createfc(c *client, namefc string, total string) {
	//create fc
	//var deckname string

	
	//c.status = "create cfc"
	//readMsg(c)
	/*
		if CheckDeckExist(db, deckname) == 0 {
			sqlStatement := `INSERT INTO Deck_instance(deckName, dateCreate) VALUES(?, NOW())`
			_, err := db.Exec(sqlStatement, deckname)
			checkErr(err)
		} else {
			fmt.Println("This Deck Name Already Used.")
		}
		var checkid int
		sqlStatement := `SELECT deckId FROM Deck_instance WHERE deckName = ? ORDER BY deckId DESC LIMIT 1 ` //check the lastest deckId and we will put it in the flashcard table
		rows, err := db.Query(sqlStatement, deckname)
		for rows.Next() {
			err = rows.Scan(&checkid)
			checkErr(err)
		}
		checkErr(err)
		fmt.Println("Number of Flashcard : ") //let user choose
		var numfc int
		fmt.Scanln(&numfc)
		var slice []cache.FlashCard
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
			checkErr(err)
		}
		cache.RedisAddDeck(redisHandler.Client, redisInstanceDeck)
	*/
}

func (s *server) login(c *client, username string, pass string) {
	//login
}

func (s *server) regis(c *client, username string, pass string, verifypass string) {
	//regis
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	msg := strings.Join(args[1:len(args)], " ")
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