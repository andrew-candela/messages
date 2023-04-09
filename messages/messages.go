package messages

import (
	"crypto/rsa"
	"fmt"

	"google.golang.org/protobuf/proto"
)

const (
	DATAGRAM_SIZE = 1023
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

func PacketFromBytes(data []byte) Packet {
	newMessage := &Message{}
	// fmt.Println("Message = ", string(data))
	err := proto.Unmarshal(data, newMessage)
	if err != nil {
		panic(err)
	}
	return Packet{
		SenderName: newMessage.SenderName,
		Content:    newMessage.Content,
		Signature:  newMessage.Signature,
		AESKey:     newMessage.AESKey,
	}
}

// If a message is too long (encoded length is over 1024 bytes)
// then it will be split into one or more datagrams
func SplitMessageIntoDatagrams(encodedPacket []byte) [][]byte {
	packetLength := len(encodedPacket)
	numGrams := (packetLength / DATAGRAM_SIZE) + 1
	expectMoreMessages := true
	var gramList [][]byte
	for i := 0; i < packetLength; i += DATAGRAM_SIZE {
		end := i + DATAGRAM_SIZE
		if end > packetLength {
			end = packetLength
		}
		if i == numGrams-1 {
			expectMoreMessages = false
		}
		newDatagramContent := encodedPacket[i:end]
		dataGram := DataGram{
			ExpectMoreMessages: expectMoreMessages,
			Content:            newDatagramContent,
		}
		encodedGram, _ := dataGram.ToBytes()

		gramList = append(gramList, encodedGram)
	}
	return gramList
}
