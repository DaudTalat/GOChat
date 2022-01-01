package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*chat_room
	requests chan request
}

// init a server
func init_server() *server {
	// returning a pointer to a server object we are making
	return &server{
		rooms:    make(map[string]*chat_room),
		requests: make(chan request),
	}
}

// init a new client
func (s *server) init_client(connection net.Conn) {
	defer connection.Close()
	// returning a client pointer to the client object we are making
	df_c := &client{
		nick:     "stranger",
		cn:       connection,
		requests: s.requests, // sets channel of client to server so that they can send request to server
	}

	log.Printf("a new client has joined the server.")

	// read message is an embedded loop that runs per each user
	df_c.read_message()
}

// join a room
func (s *server) join(c *client, args []string) {

	if len(args) == 1 {
		c.send_msg("Enter a valid room name" + "\n")
		return
	}

	if len(args[1]) <= 3 {
		c.send_msg("Please enter a valid room name" + "\n")
		return
	}

	// get command args for room name
	room_name := args[1]

	if c.room != nil {
		if c.room.name == room_name {
			c.send_msg("You are already in this room" + "\n")
			return
		}
	}

	// use the room name gotten from client to get server room
	r, exist := s.rooms[room_name]

	if !exist {
		c.send_msg("There was no instance of room: " + room_name + "\n")
		return
	}

	// remove client from previous room before joining a new one
	s.quit_room(c)

	// add client to the room members
	r.members[c.cn.RemoteAddr()] = c

	// set client to correct room
	c.room = r
	r.declare(c, fmt.Sprintf("%s has joined the room", c.nick))
}

// create a room
func (s *server) create(c *client, args []string) {

	if len(args) == 1 {
		return
	}

	if len(args[1]) <= 3 {
		print("Please enter a valid room name" + "\n")
		return
	}

	// get command args for room name
	room_name := args[1]
	// use the room name gotten from client to get server room
	// If the room doesn't exist we need to make one and set it to r
	if s.rooms[room_name] == nil {

		r := &chat_room{
			name:    room_name,
			members: make(map[net.Addr]*client),
		}

		s.rooms[room_name] = r

		print("Room '" + room_name + "' has been created" + "\n")

	} else {
		c.send_msg("There was already an instance of room: " + room_name)
		return
	}
}

// set nick to arg
func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.send_msg(fmt.Sprintf("Nickname changed to '%s'", c.nick))
}

// print all rooms
func (s *server) all_rooms(c *client, args []string) {
	for name := range s.rooms {
		c.send_msg(name)
	}

	if len(s.rooms) == 0 {
		c.send_msg("There are currently no rooms" + "\n")
	}
}

func (s *server) send_msg(c *client, args []string) {

	if len(args) == 1 || len(args[1]) == 0 {
		return
	}

	if c.room == nil {
		c.send_msg("Please join an open room")
	} else {
		msg := strings.Join(args[1:], " ")
		c.room.declare(c, c.nick+": "+msg)
	}
}

// leave room
func (s *server) quit_room(c *client) {
	// store room in temp variable before leaving, leave, then send message
	if c.room != nil {
		prev_room := s.rooms[c.room.name]
		prev_room.declare(c, fmt.Sprintf("%s has left the room", c.nick))
		delete(s.rooms[c.room.name].members, c.cn.RemoteAddr())

		c.room = nil
	}
}

// client exits
func (s *server) exit(c *client) {
	s.quit_room(c)
	c.cn.Close()
	log.Printf("a client has left the server.")
}

// case handling for server from client
func (s *server) run() {
	for cmd := range s.requests { // Go through FIFO
		switch cmd.id {
		case CMD_CREATE:
			s.create(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit_room(cmd.client)
		case CMD_ROOMS:
			s.all_rooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.send_msg(cmd.client, cmd.args)
		case CMD_EXIT:
			s.exit(cmd.client)
		}
	}
}
