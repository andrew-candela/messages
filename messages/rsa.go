package messages

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func EncryptMessage(publicKey rsa.PublicKey, message []byte, label []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		message,
		label,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return encrypted, nil
}

func DecryptMessage(privateKey rsa.PrivateKey, message []byte) []byte {
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		nil,
		&privateKey,
		message,
	)
}
