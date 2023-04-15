package messages

import (
	"crypto/rsa"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

const (
	// udpBuffer is 1024 bytes.
	// Marshalling as Protobuf adds 5 bytes.
	// We'll use 1000 bytes here to give some leeway
	DATAGRAM_SIZE = 1019
)

type DataGram struct {
	ExpectMoreMessages bool
	Content            []byte
}

func (datagram *DataGram) ToBytes() ([]byte, int) {
	g := &ProtoGram{
		ExpectMoreMessages: datagram.ExpectMoreMessages,
		Content:            datagram.Content,
	}
	data, err := proto.Marshal(g)
	if err != nil {
		panic(err)
	}
	return data, len(data)
}

func DataGramFromBytes(data []byte) (DataGram, error) {
	newGram := &ProtoGram{}
	err := proto.Unmarshal(data, newGram)
	if err != nil {
		return DataGram{}, fmt.Errorf("could not unmarshal raw bytes into DataGram...%w", err)
	}
	return DataGram{
		ExpectMoreMessages: newGram.ExpectMoreMessages,
		Content:            newGram.Content,
	}, nil
}

type Packet struct {
	SenderName string
	Content    []byte
	Signature  []byte
	AESKey     []byte
}

func (packet *Packet) ToBytes() []byte {
	m := &Message{
		SenderName: packet.SenderName,
		Content:    packet.Content,
		Signature:  packet.Signature,
		AESKey:     packet.AESKey,
	}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return data
}

func (packet *Packet) DecryptContent(key *rsa.PrivateKey) error {
	aesKey, err := RSADecrypt(key, packet.AESKey)
	if err != nil {
		return fmt.Errorf("could not decrypt AES Key...%w", err)
	}
	content, err := AESDecrypt(packet.Content, aesKey)
	if err != nil {
		return fmt.Errorf("could not decrypt packet content...%w", err)
	}
	packet.Content = []byte(content)
	return nil
}

func PacketFromBytes(data []byte) (Packet, error) {
	newMessage := &Message{}
	// fmt.Println("Message = ", string(data))
	err := proto.Unmarshal(data, newMessage)
	if err != nil {
		return Packet{}, err
	}
	return Packet{
		SenderName: newMessage.SenderName,
		Content:    newMessage.Content,
		Signature:  newMessage.Signature,
		AESKey:     newMessage.AESKey,
	}, nil
}

// If a message is too long (encoded length is over 1024 bytes)
// then it will be split into one or more datagrams
func SplitMessageIntoDatagrams(encodedPacket []byte) [][]byte {
	packetLength := len(encodedPacket)
	expectMoreMessages := true
	var gramList [][]byte
	for i := 0; i < packetLength; i += DATAGRAM_SIZE {
		end := i + DATAGRAM_SIZE
		if end > packetLength {
			end = packetLength
		}
		if end == packetLength {
			expectMoreMessages = false
		}
		newDatagramContent := encodedPacket[i:end]
		dataGram := DataGram{
			ExpectMoreMessages: expectMoreMessages,
			Content:            newDatagramContent,
		}
		encodedGram, _ := dataGram.ToBytes()
		if len(encodedGram) > 1024 {
			fmt.Println("Encoded Datagram is too long! Exiting the app")
			os.Exit(1)
		}

		gramList = append(gramList, encodedGram)
	}
	return gramList
}
