package main

import (
	"github.com/andrew-candela/messages/messages"
)

const (
	HOST = "10.0.0.186"
	PORT = "1053"
	TYPE = "udp"
)

const USERNAME = "Andrew"

func main() {
	messages.ProduceMessages(HOST+":"+PORT, USERNAME)
}
