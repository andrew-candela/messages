package messages

import (
	"os"
	"testing"
)

func TestRSAEncode(t *testing.T) {
	message := "hello world"
	k := GenerateRandomKey()
	cipher, _ := RSAEncrypt(&k.PublicKey, []byte(message))
	decoded, _ := RSADecrypt(k, cipher)
	if string(decoded) != message {
		t.Errorf("%s != %s", message, string(decoded))
	}
}

func TestRSAVerify(t *testing.T) {
	message := []byte("Hello world!")
	k := GenerateRandomKey()
	sig, _ := RSASign(k, message)
	if !RSAVerify(&k.PublicKey, message, sig) {
		t.Errorf("Verification Failed!")
	}
}

func TestRSAWriteRead(t *testing.T) {
	message := []byte("hello world")
	test_key_file := "test_key.pem"
	k := GenerateRandomKey()
	WriteKeyToDisk(k, test_key_file)
	defer os.Remove(test_key_file)
	new_key, _ := ReadExistingKey(test_key_file)
	cipher, _ := RSAEncrypt(&new_key.PublicKey, message)
	decoded, _ := RSADecrypt(new_key, cipher)
	if string(decoded) != string(message) {
		t.Errorf("%s != %s", message, decoded)
	}
}
