package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_LOGIN
	CMD_REGIS
	CMD_START
	CMD_CROOM
	CMD_CUSER
	CMD_SRD
)

type command struct {
	id     commandID
	client *client
	args   []string
}
