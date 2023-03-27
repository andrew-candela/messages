package messages

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
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
		new_err := fmt.Errorf("trouble encrypting string...%w", err)
		return nil, new_err
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
func ReadExistingKey(keyFile string) (*rsa.PrivateKey, error) {
	keyfile, err := os.ReadFile(keyFile)
	CheckErrFatal(err)
	// block, _ := pem.Decode([]byte(keyfile))
	// key_interface, err := ssh.ParseRawPrivateKey(block.Bytes)
	key, err := ssh.ParseRawPrivateKey(keyfile)
	if err != nil {
		err = fmt.Errorf("error parsing key file...%w", err)
		return nil, err
	}
	rsaKey := key.(*rsa.PrivateKey)
	return rsaKey, nil
}

func WriteKeyToDisk(key *rsa.PrivateKey, fileName string) {
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	err := os.WriteFile(fileName, pemData, 0600)
	CheckErrFatal(err)
}

func GenerateRandomKey() *rsa.PrivateKey {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	CheckErrFatal(err)
	return k
}
