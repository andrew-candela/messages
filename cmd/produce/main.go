package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/andrew-candela/messages/messages"
)

// const (
// 	HOST = "localhost"
// 	PORT = "8080"
// 	TYPE = "tcp"
// )

const USERNAME = "Andrew"

func main() {
	udpServer, err := net.ResolveUDPAddr("udp", ":1053")

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()
	// infinite loop waiting to get input
	fmt.Print("Your outbound IP is ", messages.GetOutboundIP())
	fmt.Println("Write your messages here.\nClose input with crtl+c")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("--> ")
	for scanner.Scan() {
		line := scanner.Bytes()
		newMessage := messages.Packet{
			SenderName: USERNAME,
			Content:    string(line),
		}
		bytes := newMessage.ToBytes()
		conn.Write(bytes)
		fmt.Print("--> ")
	}
}
