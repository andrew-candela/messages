package messages

import (
	"os"
	"reflect"
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

func TestBytesToString(t *testing.T) {
	sig := []byte{1, 4, 6, 8}
	expected := "[1,4,6,8]"
	got := BytesToString(sig)
	if expected != got {
		t.Errorf("%v != %v", expected, got)
	}
}

func TestStringToBytes(t *testing.T) {
	sig_str := "[1,4,6,8]"
	expected := []byte{1, 4, 6, 8}
	got, err := StringToBytes(sig_str)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("%v != %v", expected, got)
	}
}

// Creates a new key, uses it to sign something
// then serialize the signature, deserialize it,
// then try to verify it.
func TestVerifySerializedSignature(t *testing.T) {
	k := GenerateRandomKey()
	message := []byte{1, 36, 72}
	sig, err := RSASign(k, message)
	if err != nil {
		t.Error(err)
	}
	sig_str := BytesToString(sig)
	message_str := BytesToString(message)
	deserialized_sig, _ := StringToBytes(sig_str)
	deserialized_message, _ := StringToBytes(message_str)
	if !RSAVerify(&k.PublicKey, deserialized_message, deserialized_sig) {
		t.Error("Unable to verify signature!")
	}
}
