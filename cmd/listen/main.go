package main

import (
	"flag"
	"fmt"

	"github.com/andrew-candela/messages/messages"
)

const PORT = "1053"

func main() {
	var port string
	var keyfile string
	flag.StringVar(&port, "port", PORT, "Port number to listen on")
	flag.StringVar(&keyfile, "keyfile", "/Users/andrewcandela/.ssh/id_rsa", "Private keyfile to use to decrypt messages.")
	key, err := messages.ReadExistingKey(keyfile)
	if err != nil {
		fmt.Println("Could not read keyfile:", err)
		return
	}
	c := make(chan []byte, 10)
	go messages.Listen(PORT, c, *key)
	messages.PrintUDPOutput(c)
}
