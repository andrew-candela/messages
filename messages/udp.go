package messages

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	PROTOCOL   = "udp4"
	X_MARK     = "\u274C"
	CHECK_MARK = "\u2705"
)

// Get outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type GroupDetails struct {
	DestinationHostPort string
	PublicKey           *rsa.PublicKey
}

func findGroupMember(hostPort string, groupData []GroupDetails) (*GroupDetails, bool) {
	targetHost := strings.Split(hostPort, ":")[0]
	for _, detail := range groupData {
		groupHost := strings.Split(detail.DestinationHostPort, ":")[0]
		if targetHost == groupHost {
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
		userDetail, found := findGroupMember(respAddr.String(), groupDetails)
		if !found {
			fmt.Println("Could not find user located at host:", respAddr.String())
			continue
		}
		dataGram, err := DataGramFromBytes(buffer[:n])
		if err != nil {
			fmt.Println("Could not unpack datagram received from:", respAddr.String())
			continue
		}
		packet := PacketFromBytes(dataGram.Content)
		err = packet.DecryptContent(&key)
		if err != nil {
			fmt.Println("Could not decrypt message from:", respAddr.String(), err)
			continue
		}
		if !RSAVerify(userDetail.PublicKey, []byte(packet.Content), packet.Signature) {
			fmt.Println("Message received was not signed by user expected at:", userDetail.DestinationHostPort)
			continue
		}
		out_chan <- packet
		_, err = connection.WriteToUDP([]byte(CHECK_MARK), respAddr)
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
	var channels []chan Packet
	wg := sync.WaitGroup{}
	for i, detail := range deliver_details {
		fmt.Println("deliver_detail:", detail.DestinationHostPort, i)
		channels = append(channels, make(chan Packet))
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
			Content:    message_bytes,
			Signature:  sig,
		}
		for i := range deliver_details {
			channels[i] <- packet
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
// a Packet must be:
//
//	processed
//	   - a random AES key is generated and encrypted, and added to the Packet
//	   - the Packet.Content is encrypted with the AES key
//	serialized
//	  - call Packet.ToBytes()
//	chunked
//	  - the serialized bytes are chunked into batches of no more than 1023 bytes
//	  - each chunk is used to create a DataGram object
//	  - each Datagram is serialized and sent
func MessageProducer(recipient GroupDetails, user string, packet_chan <-chan Packet, wg *sync.WaitGroup) {
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
	for packet := range packet_chan {
		if !made_connection {
			fmt.Println(recipient.DestinationHostPort, X_MARK)
			wg.Done()
			continue
		}

		preparedPacket, err := processPacket(&packet, recipient.PublicKey)
		if err != nil {
			fmt.Println("Could not process Packet...%w", err)
			wg.Done()
			continue
		}
		grams := SplitMessageIntoDatagrams(preparedPacket)
		for _, gram := range grams {
			resp_buffer := make([]byte, 1024)
			_, err = conn.Write(gram)
			if err != nil {
				fmt.Println("Failed to send message to", recipient.DestinationHostPort, err)
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
}

// Prepares a Packet for sending.
// Creates a new AES key, encrypts the Content, encrypts the key
// and finally serializes the packet as bytes
func processPacket(packet *Packet, pubKey *rsa.PublicKey) ([]byte, error) {
	aesKey, err := GenerateRandomAESKey()
	if err != nil {
		return nil, fmt.Errorf("could not generate AES key...%w", err)
	}
	encryptedKey, err := RSAEncrypt(pubKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("could not RSA Encrypt AES key.. %w", err)
	}
	encryptedContent, err := AESEncrypt(string(packet.Content), aesKey)
	if err != nil {
		return nil, fmt.Errorf("could not AES Encrypt message..%w", err)
	}
	packet.Content = encryptedContent
	packet.AESKey = encryptedKey
	return packet.ToBytes(), nil

}
