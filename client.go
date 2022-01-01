package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// client type has a nick, designated chat roo, cn I/O, request header, and a channel (one value buffer between concurrent goroutines)
type client struct {
	nick     string
	room     *chat_room
	cn       net.Conn
	requests chan<- request // write only channel
}

// read messages, pair type of request,
func (c *client) read_message() {
	for {
		msg, err := bufio.NewReader(c.cn).ReadString('\n') // A delim term, signifies end of message.

		if err != nil {
			return
		} else {
			msg = strings.Trim(msg, "\r\n")
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])

			// request builder
			switch cmd {
			case "/create":
				c.requests <- request{
					id:     CMD_CREATE,
					client: c,
					args:   args,
				}
			case "/join":
				c.requests <- request{
					id:     CMD_JOIN,
					client: c,
					args:   args,
				}
			case "/nick":
				c.requests <- request{
					id:     CMD_NICK,
					client: c,
					args:   args,
				}
			case "/quit":
				c.requests <- request{
					id:     CMD_QUIT,
					client: c,
				}
			case "/rooms":
				c.requests <- request{
					id:     CMD_ROOMS,
					client: c,
				}
			case "/msg":
				c.requests <- request{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			case "/exit":
				c.requests <- request{
					id:     CMD_EXIT,
					client: c,
				}

			default:
				c.send_err(fmt.Errorf("invalid command: %s", cmd))
			}
		}
	}
}

func (c *client) send_err(err error) {
	c.cn.Write([]byte("ERROR: " + err.Error() + "\n"))
}

func (c *client) send_msg(msg string) {
	c.cn.Write([]byte("> " + msg + "\n"))
}
