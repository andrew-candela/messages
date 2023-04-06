package messages

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	PROTOCOL = "udp4"
	X_MARK   = "\u274C"
)

type GroupDetails struct {
	DestinationHostPort string
	PublicKey           rsa.PublicKey
}

func findGroupMember(hostPort string, groupData []GroupDetails) (*GroupDetails, bool) {
	for _, detail := range groupData {
		if detail.DestinationHostPort == hostPort {
			return &detail, true
		}
	}
	return nil, false
}

// Creates a UDP "server", listening on the given port
// Will loop forever and pass input to the given channel
func Listen(port string, out_chan chan Packet, key rsa.PrivateKey, groupDetails []GroupDetails) {
	udpAddr, err := net.ResolveUDPAddr(PROTOCOL, ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening at: ", GetOutboundIP().String()+udpAddr.String())
	connection, err := net.ListenUDP(PROTOCOL, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	for {
		buffer := make([]byte, 1024)
		n, respAddr, _ := connection.ReadFromUDP(buffer)
		message := RSADecrypt(key, buffer[:n])
		parsed := PacketFromBytes(message)
		userDetail, found := findGroupMember(respAddr.String(), groupDetails)
		if !found {
			fmt.Println("Could not find user located at host:", respAddr.String())
			continue
		}
		if !RSAVerify(&userDetail.PublicKey, []byte(parsed.Content), parsed.Signature) {
			fmt.Println("Message received was not signed by user expected at:", userDetail.DestinationHostPort)
			continue
		}
		out_chan <- parsed
		_, err := connection.WriteToUDP([]byte("\u2705"), respAddr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Prints the data received from the input channel.
// Pair this with Listen to print messages you got from friends.
func PrintUDPOutput(in_channel <-chan Packet) {
	for message := range in_channel {
		fmt.Printf("%s: %s\n", message.SenderName, message.Content)
	}
}

// Write UDP messages to a given IP:PORT address
func ProduceMessages(deliver_details []GroupDetails, user string, key *rsa.PrivateKey) {
	var channels []chan []byte
	wg := sync.WaitGroup{}
	for i, detail := range deliver_details {
		fmt.Println("deliver_detail:", detail.DestinationHostPort, i)
		channels = append(channels, make(chan []byte))
		go MessageProducer(detail, user, channels[i], &wg)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Println("Write your messages below. Close input with crtl+c:")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("--> ")
	for scanner.Scan() {
		message_bytes := scanner.Bytes()
		sig, err := RSASign(key, message_bytes)
		if err != nil {
			fmt.Println("Error signing message. Will not send!")
			continue
		}
		packet := Packet{
			SenderName: user,
			Content:    string(message_bytes),
			Signature:  sig,
		}
		proto_message := packet.ToBytes()
		for i := range deliver_details {
			channels[i] <- proto_message
			wg.Add(1)
		}
		wg.Wait()
		fmt.Print("--> ")
	}

}

func createUDPConnection(ip_port string) (conn *net.UDPConn, err error) {
	udpAddr, err := net.ResolveUDPAddr(PROTOCOL, ip_port)
	if err != nil {
		return nil, err
	}
	conn, err = net.DialUDP(PROTOCOL, nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Write UDP messages to a given IP:PORT address
func MessageProducer(recipient GroupDetails, user string, message_chan <-chan []byte, wg *sync.WaitGroup) {
	made_connection := true
	conn, err := createUDPConnection(recipient.DestinationHostPort)
	if err != nil {
		made_connection = false
		fmt.Println("Could not make connection to", recipient.DestinationHostPort, "...", err)
	}
	// signal that we are done establishing the connection
	wg.Done()

	// When a message is published to the channel, the calling function blocks
	// until we signal that the message has been processed.
	for message := range message_chan {
		if !made_connection {
			fmt.Println(recipient.DestinationHostPort, X_MARK)
			wg.Done()
			continue
		}
		encrypted_message, err := RSAEncrypt(recipient.PublicKey, message)
		if err != nil {
			fmt.Println("Failed to send message to", recipient.DestinationHostPort, err)
			wg.Done()
			continue
		}
		resp_buffer := make([]byte, 1024)
		_, err = conn.Write(encrypted_message)
		if err != nil {
			fmt.Println(recipient.DestinationHostPort, X_MARK)
			wg.Done()
			continue
		}
		n, _, err := conn.ReadFromUDP(resp_buffer)
		if err == nil {
			fmt.Println(recipient.DestinationHostPort, string(resp_buffer[:n]))
		} else {
			fmt.Println(recipient.DestinationHostPort, X_MARK)
		}
		wg.Done()
	}
}
