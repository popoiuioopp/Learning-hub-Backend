package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	status   string
	commands chan<- command
	score    int
	vaild    bool
}

func (c *client) readInput() {
	for {
		if c.status == "0" {
			c.msg(fmt.Sprint("----Ready for your command----"))
			// ip := c.conn.RemoteAddr().String()
			// c.msg(fmt.Sprintf("Test Print %s", ip))
			msg, err := bufio.NewReader(c.conn).ReadString('\n')
			if err != nil {
				c.msg(fmt.Sprintf(msg))
				return
			}

			msg = strings.Trim(msg, "\r\n")

			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])

			switch cmd {
			case "/nick":
				c.commands <- command{
					id:     CMD_NICK,
					client: c,
					args:   args,
				}
			case "/join":
				c.commands <- command{
					id:     CMD_JOIN,
					client: c,
					args:   args,
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
				c.commands <- command{
					id:     CMD_CROOM,
					client: c,
					args:   args,
				}
			case "/cfc":
				c.status = "1"
			case "/rfc":
				c.status = "2"
			case "/rstatus":
				c.msg(fmt.Sprintf("room status: %t ,current deckid:%d,current host_id:%s\n", c.room.status, c.room.deck.deckID, c.room.host))
			case "/ready":
				c.status = "3"
				c.room.Changeroomstatus(c)
			case "/srd":
				c.commands <- command{
					id:     CMD_SRD,
					client: c,
					args:   args,
				}
			case "/cuser":
				c.msg(fmt.Sprintf("current user: "))
			default:
				c.err(fmt.Errorf("unknown command: %s", cmd))
			}
		} else if c.status == "1" {
			c.msg(fmt.Sprintf("Name Your Deck: "))
			deckname, err := bufio.NewReader(c.conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			if !createDeck(c, deckname) {
				c.status = "0"
				c.msg(fmt.Sprintf("you are in the lobby now "))
			}
			deckid, err := checkDeckId(deckname)
			if err != nil {
				fmt.Println(err)
				return
			}
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
						c.msg(fmt.Sprintf("You are in lobby now"))

						// REMOVE IT LATER ON
						// msg123 := "/cfc text 123"
						// args := strings.Split(msg123, " ")

						// c.commands <- command{
						// 	id: CMD_CREATEFC,
						// 	client: c,
						// 	args:	args,
						// }
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
					tempFC.deckID = deckid
					tempFC.definition = def
					tempFC.term = term
					fcList = append(fcList, tempFC)
				} else {
					c.msg(fmt.Sprintf("Invalid inputs"))
					continue
				}
			}
		} else if c.status == "2" {
			ListDecks(c)
			deckid, err := bufio.NewReader(c.conn).ReadString('\n')
			if err != nil {
				c.msg(fmt.Sprintf(deckid))
				return
			}
			c.msg(fmt.Sprintf("Deckid: %s\n", deckid))
			c.msg(fmt.Sprintf("Pls Choose your deck by type in deck id"))
			//query fc from db

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

			c.status = "0"

		} else if c.status == "3" {
			for c.room.status == true {
				msg, err := bufio.NewReader(c.conn).ReadString('\n')
				msg = strings.Trim(msg, "\r\n")
				if err != nil {
					c.msg(fmt.Sprintf(msg))
					return
				}
				if c.vaild {
					if msg == c.room.answer {
						c.msg(fmt.Sprintf("Correct!"))
						c.score += 1
					} else {
						c.msg(fmt.Sprintf("Try Again!"))
					}
				} else {
					c.msg(fmt.Sprintf("You Already Answer!"))
				}

				c.vaild = false
			}

		} else if c.status == "boardcast" {

			if c.room.status == true {
				c.room.GenQuestion(c)
			}

		}
	}

}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
