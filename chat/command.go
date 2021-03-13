package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_CROOM
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG 
	CMD_QUIT
	CMD_LOGIN
	CMD_REGIS
	CMD_CREATEFC
	CMD_START
)

type command struct {
	id     commandID
	client *client
	args   []string
}
