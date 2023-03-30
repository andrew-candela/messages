package messages

import (
	"os"
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

func TestRSAWriteRead(t *testing.T) {
	message := []byte("hello world")
	test_key_file := "test_key.pem"
	k := GenerateRandomKey()
	WriteKeyToDisk(k, test_key_file)
	defer os.Remove(test_key_file)
	new_key, _ := ReadExistingKey(test_key_file)
	cipher, _ := RSAEncrypt(new_key.PublicKey, message)
	decoded := RSADecrypt(*new_key, cipher)
	if string(decoded) != string(message) {
		t.Errorf("%s != %s", message, decoded)
	}
}
