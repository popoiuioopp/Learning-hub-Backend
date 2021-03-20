package main

import (
	"fmt"
	"net"
	"time"
)

type room struct {
	name    string
	members map[net.Addr]*client
	currFC  Flashcard
	deck    Deck
	host    string
	status  bool
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			// m.status = "4"
			// m.msg(fmt.Sprintf("Client %s----> %s", addr, m.status))
			m.msg(msg)
		}
	}
}

func (r *room) changeroomstatus(sender *client) {
	for _, m := range r.members {
		if m.status != "3" {
			// sender.msg(fmt.Sprintf("Client [%s] -> %s", addr, m.status))
			r.broadcast(m, fmt.Sprintf("These users are nor ready: %s\n", m.conn.RemoteAddr().String()))
			sender.msg(fmt.Sprintf("Cannot Start!"))
			return
		}
	}
	sender.msg(fmt.Sprintf("Game Start!"))
	r.status = true
	return
}

func (r *room) GenQuestion(sender *client) {
	for _, fc := range *r.deck.fcArray {
		for addr, m := range r.members {
			if sender.conn.RemoteAddr() == addr {
				m.msg(fmt.Sprintf("%s\n", fc.definition))
				r.currFC = fc
			}
			time.Sleep(5 * time.Second)
		}
	}
	r.status = false
	for _, m := range r.members {
		m.status = "0"
	}
}
