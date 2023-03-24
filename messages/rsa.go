package messages

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const LABEL = "myCoolMessagingApp"

func RSAEncrypt(publicKey rsa.PublicKey, message []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&publicKey,
		[]byte(message),
		[]byte(LABEL),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return encrypted, nil
}

func RSADecrypt(privateKey rsa.PrivateKey, message []byte) []byte {
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		nil,
		&privateKey,
		message,
		[]byte(LABEL),
	)
	CheckErrFatal(err)
	return decrypted
}

// // Reads an existing .pem or rsa keyfile and returns a
// // reference to it.
// func ReadExistingKey(keyFile string) *rsa.PrivateKey {
// 	keyfile, err := os.ReadFile(keyFile)
// 	CheckErrFatal(err)
// 	block, _ := pem.Decode([]byte(keyfile))
// 	switch block.Type {
// 	case "OPENSSH PRIVATE KEY":
// 		key, err := ssh.ParsePrivateKey(block.Bytes)
// 		CheckErrFatal(err)
// 		return &key.PublicKey()
// 	default:
// 		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 		CheckErrFatal(err)
// 		return key
// 	}
// }

func WriteKeyToDisk(key *rsa.PrivateKey, fileName string) {
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA Private Key",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	err := os.WriteFile(fileName, pemData, 0600)
	CheckErrFatal(err)
}

func GenerateRandomKey() *rsa.PrivateKey {
	k, err := rsa.GenerateKey(rand.Reader, 1024)
	CheckErrFatal(err)
	return k
}
