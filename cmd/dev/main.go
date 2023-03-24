package main

import (
	"fmt"

	"github.com/andrew-candela/messages/messages"
)

func main() {
	// k := messages.GenerateRandomKey()
	// messages.WriteKeyToDisk(k, "test_key.pem")
	k := messages.ReadExistingKey("/Users/andrewcandela/.ssh/id_rsa")
	message := "hey there!"
	cipher, _ := messages.RSAEncrypt(k.PublicKey, []byte(message))
	decoded := messages.RSADecrypt(*k, cipher)
	fmt.Println(string([]byte(decoded)))
}
