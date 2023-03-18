package main

import (
	"fmt"
	"log"
	"net"

	"github.com/andrew-candela/messages/messages"
)

func main() {
	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	for {
		buf := make([]byte, 1024)
		n, _, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		message := messages.FromBytes(buf[:n])
		fmt.Printf("%s: %s\n", message.SenderName, message.Content)
	}

}
