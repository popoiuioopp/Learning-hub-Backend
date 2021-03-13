package main

import (
	"net"
)

type room struct {
	name      string
	members   map[net.Addr]*client
	flashcard string
	host      string
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}
