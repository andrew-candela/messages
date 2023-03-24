package messages

import (
	"testing"
)

func TestCryptoEncode(t *testing.T) {
	message := "hello world"
	k := GenerateRandomKey()
	cipher, _ := RSAEncrypt(k.PublicKey, []byte(message))
	decoded := RSADecrypt(*k, cipher)
	if string(decoded) != message {
		t.Errorf("%s != %s", message, string(decoded))
	}
}
