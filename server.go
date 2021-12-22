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

	if len(args[1]) <= 3 {
		print("Please enter a valid room name.")
		return
	}

	// get command args for room name
	room_name := args[1]
	// use the room name gotten from client to get server room
	r, exist := s.rooms[room_name]

	// If the room doesn't exist we need to make one and set it to r
	if !exist {
		r = &chat_room{
			name:    room_name,
			members: make(map[net.Addr]*client),
		}

		s.rooms[room_name] = r
	}

	// add client to the map
	r.members[c.cn.RemoteAddr()] = c

	// remove client from previous room before joining a new one
	s.leave_room(c)

	// set client to correct room
	c.room = r
	r.declare(c, fmt.Sprintf("%s has joined the room", c.nick))
}

// set nick to arg
func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.send_msg(fmt.Sprintf("Nickname changed to '%s'", c.nick))
}

// print all rooms
func (s *server) all_rooms(c *client, args []string) {
	for name := range s.rooms {
		fmt.Println(name)
	}
}

func (s *server) send_msg(c *client, args []string) {
	if c.room == nil {
		c.send_msg("Please join an open room")
	} else {
		msg := strings.Join(args[1:], " ")
		c.room.declare(c, c.nick+": "+msg)
	}
}

// leave room
func (s *server) leave_room(c *client) {
	// store room in temp variable before leaving, leave, then send message
	if c.room != nil {
		prev_room := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.cn.RemoteAddr())

		prev_room.declare(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

// case handling for server from client
func (s *server) run() {
	for cmd := range s.requests {
		switch cmd.id {
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_QUIT:
			s.leave_room(cmd.client)
		case CMD_ROOMS:
			s.all_rooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.send_msg(cmd.client, cmd.args)
		}
	}
}
