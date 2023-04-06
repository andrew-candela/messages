package messages

import (
	"log"
	"net"

	"google.golang.org/protobuf/proto"
)

type Packet struct {
	SenderName string
	Content    string
	Signature  []byte
}

func (packet *Packet) ToBytes() []byte {
	m := &Message{
		SenderName: packet.SenderName,
		Content:    packet.Content,
		Signature:  packet.Signature,
	}
	data, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}
	return data
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
	}

}

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
