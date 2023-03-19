package messages

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const PROTOCOL = "udp4"

// Creates a UDP "seerver", listening on the given port
// Will loop forever and pass input to the given channel
func Listen(port string, out_chan chan []byte) {
	udpAddr, err := net.ResolveUDPAddr(PROTOCOL, ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening at: ", GetOutboundIP().String()+":"+udpAddr.String())
	connection, err := net.ListenUDP(PROTOCOL, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		n, respAddr, _ := connection.ReadFromUDP(buffer)
		out_chan <- buffer[:n]
		_, err := connection.WriteToUDP([]byte{1}, respAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Write UDP messages to a given IP:PORT address
func ProduceMessages(ip_port string, user string) {
	udpAddr, err := net.ResolveUDPAddr(PROTOCOL, ip_port)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP(PROTOCOL, nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The UDP server is %s\n", conn.RemoteAddr().String())
	defer conn.Close()
	fmt.Println("Write your messages here.\nClose input with crtl+c")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("--> ")
	for scanner.Scan() {
		resp_buffer := make([]byte, 1024)
		line := scanner.Bytes()
		newMessage := Packet{
			SenderName: user,
			Content:    string(line),
		}
		bytes := newMessage.ToBytes()
		conn.Write(bytes)
		n, _, err := conn.ReadFromUDP(resp_buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Reply: %s\n", string(resp_buffer[:n]))
		fmt.Print("--> ")
	}

}

// Prints the data received form the input channel.
// Pair this with Listen to print messages you got from friends.
func PrintUDPOutput(in_channel chan []byte) {
	for datagram := range in_channel {
		message := FromBytes(datagram)
		fmt.Printf("%s: %s\n", message.SenderName, message.Content)
	}
}
