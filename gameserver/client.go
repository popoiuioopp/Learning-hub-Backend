package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/go-redis/redis"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	status   string
	commands chan<- command
	score    int
	vaild    bool
	no_ques  int
}

func (c *client) readInput() {

	// read room from user
	roomname, _ := bufio.NewReader(c.conn).ReadString('\n')
	arr := []string{" ", roomname}

	c.commands <- command{
		id:     CMD_JOIN,
		client: c,
		args:   arr,
	}

	for {

		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			c.msg(fmt.Sprintf(msg))
			return
		}
		msg = strings.Trim(msg, "\r\n")
		c.msg(fmt.Sprint(msg))

		if c.status == "0" {
			// c.msg(fmt.Sprint("----Ready for your command----"))
			// ip := c.conn.RemoteAddr().String()
			// c.msg(fmt.Sprintf("Test Print %s", ip))
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])

			switch cmd {
			case "/nick":
				if len(args) == 2 {
					c.commands <- command{
						id:     CMD_NICK,
						client: c,
						args:   args,
					}
				} else {
					c.msg(fmt.Sprintf("Invalid syntax"))
				}
			case "/join":
				if len(args) == 2 {
					c.commands <- command{
						id:     CMD_JOIN,
						client: c,
						args:   args,
					}
				} else {
					c.msg(fmt.Sprintf("Invalid syntax"))
				}
			case "/rooms":
				c.commands <- command{
					id:     CMD_ROOMS,
					client: c,
				}
			case "/msg":
				c.commands <- command{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			case "/quit":
				c.commands <- command{
					id:     CMD_QUIT,
					client: c,
				}
			case "/croom":
				if len(args) == 2 {
					c.commands <- command{
						id:     CMD_CROOM,
						client: c,
						args:   args,
					}
				} else {
					c.msg(fmt.Sprintf("Invalid syntax"))
				}
			case "/cfc":
				c.status = "1"
				c.msg(fmt.Sprintf("Name Your Deck: "))
			case "/rfc":
				c.status = "2"
				ListDecks(c)
				c.msg(fmt.Sprintf("Pls Choose your deck by type in deck id"))
			case "/rstatus":
				if c.room != nil {
					c.msg(fmt.Sprintf("room status: %t ,current deckid:%d,current host_id:%s\n",
						c.room.status, c.room.deck.DeckID, c.room.host))
				} else {
					c.msg(fmt.Sprintf("Please Select Room first"))
				}
			case "/ready":
				if c.room.host == c.conn.RemoteAddr().String() {
					if c.room.deck.DeckID == 0 {
						c.msg(fmt.Sprintf("Please Specify Your Deck First"))
					} else {
						c.status = "broadcast"
					}
				} else {
					c.status = "3"
				}
				c.room.Changeroomstatus(c)
			case "/srd":
				if len(args) == 2 {
					c.commands <- command{
						id:     CMD_SRD,
						client: c,
						args:   args,
					}
				} else {
					c.msg(fmt.Sprintf("Invalid syntax"))
				}
			case "/cuser":
				c.msg(fmt.Sprintf("current user: "))
			default:
				c.err(fmt.Errorf("unknown command: %s", cmd))
			}
		} else if c.status == "1" {

			deckname := msg
			ok, deckid := createDeck(c, deckname)
			if !ok {
				c.status = "0"
				c.msg(fmt.Sprintf("you are in the lobby now "))
			}
			// deckid, err := checkDeckId(deckname)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }
			c.msg(fmt.Sprintf("You deck ID is %d", deckid))
			var fcList []Flashcard

			for c.status == "1" {
				c.msg(fmt.Sprintf("write your words then space bar and || with definiton "))
				c.msg(fmt.Sprintf("Give Your Word (format: Term || Definition) or Exit (cmd: /done): "))
				msg, err := bufio.NewReader(c.conn).ReadString('\n')
				if err != nil {
					c.msg(fmt.Sprintf(msg))
					return
				}

				text := strings.Split(msg, " || ")

				if len(text) == 1 {
					cmd := strings.TrimSpace(text[0])
					switch cmd {
					case "/done":
						c.status = "0"
						createfc(c, fcList)
						fcList = nil
						c.msg(fmt.Sprintf("You are in lobby now"))
						c.msg(fmt.Sprintf("Done creating flashcard"))
						break
					default:
						c.msg(fmt.Sprintf("Invalid inputs"))
						continue
					}
				} else if len(text) == 2 {
					var tempFC Flashcard
					term := strings.TrimSpace(text[0])
					def := strings.TrimSpace(text[1])
					c.msg(fmt.Sprintf(term))
					c.msg(fmt.Sprintf("%s\n", def))
					tempFC.DeckID = deckid
					tempFC.Definition = def
					tempFC.Term = term
					fcList = append(fcList, tempFC)
				} else {
					c.msg(fmt.Sprintf("Invalid inputs"))
					continue
				}
			}
		} else if c.status == "2" {

			deckid := msg
			c.msg(fmt.Sprintf("Deckid: %s\n", deckid))

			unmarServer, err := redisServerHandler.Client.Get(deckid).Result()
			if err == redis.Nil {
				log.Printf("Can not find deck:%s in server's redis\n", deckid)
				unmar, err := redisHandler.Client.Get(deckid).Result()
				if err == redis.Nil {
					// query
					statement := "SELECT term , definition FROM Flashcard_instance WHERE deckId = ?;"
					rows, err := sqliteHandler.Conn.Query(statement, deckid)
					if err != nil {
						fmt.Print(err)
					}

					for rows.Next() {
						var term string
						var definition string
						err = rows.Scan(&term, &definition)
						if err != nil {
							fmt.Print(err)
						}
						c.msg(fmt.Sprintf("%s : %s\n", term, definition))
					}

					rows.Close() //good habit to close
				} else {
					b := []byte(unmar)
					deck := &Deck{}
					err = json.Unmarshal(b, deck)
					if err != nil {
						return
					}
					for _, item := range *deck.FcArray {
						c.msg(fmt.Sprintf("%s : %s\n", item.Term, item.Definition))
					}
				}
			} else {
				log.Printf("Pulling deck:%s from server's redis\n", deckid)
				b := []byte(unmarServer)
				deck := &Deck{}
				err = json.Unmarshal(b, deck)
				if err != nil {
					return
				}
				c.room.deck = *deck
			}

			c.status = "0"

		} else if c.status == "3" {

			if c.room.status == true {
				if c.no_ques <= c.room.no_fc {
					if c.vaild {
						if msg == c.room.answer {
							c.msg(fmt.Sprintf("Correct!"))
							c.score += 1
							if c.no_ques == c.room.no_fc {
								c.no_ques += 1
							}
						} else {
							c.msg(fmt.Sprintf("Try Again!"))
							if c.no_ques == c.room.no_fc {
								c.no_ques += 1
							}
						}
						c.vaild = false
					} else {
						c.msg(fmt.Sprintf("You Already Answer!"))
					}

				} else {
					c.msg(fmt.Sprintf("Game finish!!!"))
					c.msg(fmt.Sprintf("Wait others to finish..."))
				}

			} else {
				c.msg(fmt.Sprintf("Wait other to be ready!"))
			}

		}

		if c.status == "broadcast" {
			c.room.GenQuestion(c)
		}
	}

}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
