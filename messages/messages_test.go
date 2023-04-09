package messages

import (
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
	pack := Packet{
		SenderName: "Andrew",
		Content:    []byte("Hello"),
		Signature:  []byte("hi"),
		AESKey:     []byte("this is a secret key"),
	}
	pack_bytes := pack.ToBytes()
	new_packet := PacketFromBytes(pack_bytes)
	if string(new_packet.Signature) != string(pack.Signature) {
		t.Errorf("Oops!")
	}
}

func TestMessageLength(t *testing.T) {
	key := GenerateRandomKey()
	aesKey, _ := GenerateRandomAESKey()
	pub := key.PublicKey
	encryptedKey, _ := RSAEncrypt(&pub, aesKey)
	sig, _ := RSASign(key, []byte("hi"))
	message := "Hello this is a short message! Hello this is a short message!"
	encrypted_content, _ := AESEncrypt(message, aesKey)
	pack := Packet{
		SenderName: "Andrew",
		Content:    encrypted_content,
		Signature:  sig,
		AESKey:     encryptedKey,
	}
	pack_bytes := pack.ToBytes()
	fmt.Println(len(pack_bytes))
	fmt.Println(PacketFromBytes(pack_bytes))
	if len(pack_bytes) >= 1024 {
		t.Errorf("Your encoded message is too long")
	}
}
