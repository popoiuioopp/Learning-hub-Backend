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
	answer  string
}

type detail struct {
	name  string
	score int
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

func (r *room) Changeroomstatus(sender *client) {
	for _, m := range r.members {
		if m.status != "3" {
			// sender.msg(fmt.Sprintf("Client [%s] -> %s", addr, m.status))
			r.broadcast(m, fmt.Sprintf("These users are nor ready: %s\n", m.conn.RemoteAddr().String()))
			sender.msg(fmt.Sprintf("Cannot Start!"))
			return
		}
	}
	sender.msg(fmt.Sprintf("Game Start!"))
	sender.status = "boardcast"
	r.status = true
	return
}

// func (r *room) CheckAnswer(sender *client, answer chan string) {
// 	for {
// 		select {
// 		case text := <-answer:
// 			fmt.Println(text)
// 		case <-time.After(3 * time.Second):
// 			fmt.Println("quit")
// 			return
// 		}
// 	}
// }

func (r *room) GenQuestion(sender *client) {
	for _, fc := range *r.deck.fcArray {
		for addr, m := range r.members {
			m.vaild = true
			r.answer = fc.term
			r.currFC = fc
			if sender.conn.RemoteAddr() != addr {
				m.msg(fmt.Sprintf("%s\n", fc.definition))
			}
		}
		time.Sleep(15 * time.Second)
	}
	r.status = false

	var name []string
	maximum := -1

	for _, m := range r.members {
		if m.score >= maximum {
			if m.score > maximum {
				maximum = m.score
				name = nil
				name = append(name, m.conn.RemoteAddr().String())
			} else {
				name = append(name, m.conn.RemoteAddr().String())
			}
		}
	}

	for _, m := range r.members {
		m.msg(fmt.Sprintf("Winner:"))
		for _, winner := range name {
			m.msg(fmt.Sprintf("%s", winner))
		}
		m.msg(fmt.Sprintf("Score:\n%d point(s)", maximum))
	}

	for _, m := range r.members {
		m.vaild = false
		m.status = "0"
	}
}

/*
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
*/
