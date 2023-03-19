package main

import (
	"github.com/andrew-candela/messages/messages"
)

const PORT = "1053"

func main() {
	c := make(chan []byte, 10)
	go messages.Listen(PORT, c)
	messages.PrintUDPOutput(c)
}
