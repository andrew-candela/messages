package messages

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// generates a random 32 byte key
func GenerateRandomAESKey() (key []byte, err error) {
	key = make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		err = fmt.Errorf("error generating key...%w", err)
		return nil, err
	}
	return key, nil
}

// encrypts a plain text message with a given key
func AESEncrypt(plainText string, key []byte) ([]byte, error) {
	text := []byte(plainText)
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("unable to create new AES Cipher...%w", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("unable to create new GCM Cipher...%w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("unable to populate the Nonce...%w", err)
	}
	return gcm.Seal(nonce, nonce, text, nil), nil
}

func AESDecrypt(cipherText []byte, key []byte) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("unable to create new AES Cipher...%w", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("unable to create new GCM Cipher...%w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("nonce size is greater than length of cipherText...%w", err)
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	decryptedBytes, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("unable to decrypt plaintext...%w", err)
	}
	return string(decryptedBytes), nil
}
