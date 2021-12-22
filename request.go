package main

type requestID int

const (
	CMD_JOIN requestID = iota
	CMD_NICK
	CMD_QUIT
	CMD_ROOMS
	CMD_MSG
)

type request struct {
	id     requestID
	client *client
	args   []string
}
