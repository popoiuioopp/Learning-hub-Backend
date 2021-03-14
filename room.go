package main

import (
	"net"
	"fmt"
)

type room struct {
	name      string
	members   map[net.Addr]*client
	flashcard string
	host      string
	status    bool
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

func (r *room) changeroomstatus(sender *client) int {
	for _ , m := range r.members {
		if (m.status != "3")  {
			// sender.msg(fmt.Sprintf("Client [%s] -> %s", addr, m.status))
			sender.msg(fmt.Sprintf("Cannot Start!"))
			return 0
		}
	}
	sender.msg(fmt.Sprintf("Game Start!"))
	r.status = true
	return 1
}