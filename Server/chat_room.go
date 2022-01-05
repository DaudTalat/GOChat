package main

import (
	"net"
)

type chat_room struct {
	name    string
	members map[net.Addr]*client // member = [Key = Address, value = client}
}

func (cr *chat_room) declare(sending_c *client, msg string) {
	// For each connection address, declare m to be members of chat room, and cn_addr to be the Address
	// Getting map data for members with the respective chat room then sending the data to all members
	for cn_addr, m := range cr.members {
		if cn_addr != sending_c.cn.RemoteAddr() { // We want to make sure we don't send a message to ourselves
			m.send_msg(msg)
		}
	}
}
