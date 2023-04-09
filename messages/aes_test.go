package messages

import (
	"fmt"
	"strings"
	"testing"
)

func TestGenerateRandomAESKey(t *testing.T) {
	k, err := GenerateRandomAESKey()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(k)
}

func TestAESEncrypt(t *testing.T) {
	plainText := "Hello Andrew!"
	key, _ := GenerateRandomAESKey()
	cipherText, _ := AESEncrypt(plainText, key)
	recoveredText, _ := AESDecrypt(cipherText, key)
	if plainText != recoveredText {
		t.Error("unexpected results!")
	}
}

func TestAESEncryptLong(t *testing.T) {
	plainText := strings.Repeat("Hello Andrew! ", 50)
	key, _ := GenerateRandomAESKey()
	cipherText, _ := AESEncrypt(plainText, key)
	recoveredText, _ := AESDecrypt(cipherText, key)
	if plainText != recoveredText {
		t.Error("unexpected results!")
	}
}
