package main

type requestID int

const (
	CMD_CREATE requestID = iota
	CMD_JOIN
	CMD_NICK
	CMD_QUIT
	CMD_ROOMS
	CMD_MSG
	CMD_EXIT
)

type request struct {
	id     requestID
	client *client
	args   []string
}
