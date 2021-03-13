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
	status 	 string
	commands chan<- command
}

func (c *client) readInput() {
	for {
		c.msg(fmt.Sprintf("Status Client1: %s\n", c.status))
		if (c.status == "0") {
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
			case "/start":
				c.commands <- command{
					id:     CMD_START,
					client: c,
					args:   args,
				}
			case "/cfc":
				c.status = "1"
			
			default:
				c.err(fmt.Errorf("unknown command: %s", cmd))
			}
		} else if (c.status == "1") {
			c.msg(fmt.Sprintf("Name Your Deck: "))
			deckname, err := bufio.NewReader(c.conn).ReadString('\n')
			if err != nil {
				return
			}
			c.msg(fmt.Sprintf("Deckname: %s\n",deckname))

			//for loop
			
			for (c.status == "1"){
				msg, err := bufio.NewReader(c.conn).ReadString('\n')
				if err != nil {
					c.msg(fmt.Sprintf(msg))
					return 
				}

				text := strings.Split(msg, " || ")
				
				if (len(text)==1){
					cmd := strings.TrimSpace(text[0])
					switch cmd {
						case "/exit":
							c.status = "0"	
							msg123 := "/cfc text 123"
							args := strings.Split(msg123, " ")

							c.commands <- command{
								id: CMD_CREATEFC,
								client: c,
								args:	args,
							} 
							break
						default:
							c.msg(fmt.Sprintf("Invalid inputs"))
					}
				} else if (len(text) == 2) {
					term := strings.TrimSpace(text[0])
					def := strings.TrimSpace(text[1])
					c.msg(fmt.Sprintf(term))
					c.msg(fmt.Sprintf(def))
					c.msg(fmt.Sprintf("Term & def"))
				} else {
					c.msg(fmt.Sprintf("Invalid inputs"))
					continue
					}		
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
