package messages

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

const LABEL = "myCoolMessagingApp"

// Convert a byte slice to a string by converting it to a slice of
// ints and serializing that to a JSON string.
func BytesToString(sig []byte) string {
	var sig_int []int
	for _, byte_val := range sig {
		sig_int = append(sig_int, int(byte_val))
	}
	sig_str, err := json.Marshal(sig_int)
	if err != nil {
		panic(err)
	}
	return string(sig_str)
}

// Take a JSON string (an array of ints) and build a []byte
func StringToBytes(sig_str string) ([]byte, error) {
	var sig_byte []byte
	var data []byte
	err := json.Unmarshal([]byte(sig_str), &data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, byte_int := range data {
		sig_byte = append(sig_byte, byte(byte_int))
	}
	return sig_byte, nil
}

func RSAEncrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	encrypted, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		[]byte(LABEL),
	)
	if err != nil {
		new_err := fmt.Errorf("trouble encrypting string...%w", err)
		return nil, new_err
	}
	return encrypted, nil
}

func RSADecrypt(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptOAEP(
		sha256.New(),
		nil,
		privateKey,
		message,
		[]byte(LABEL),
	)
	if err != nil {
		new_err := fmt.Errorf("unable to decrypt message... %w", err)
		return nil, new_err
	}
	return decrypted, nil
}

func RSASign(key *rsa.PrivateKey, message []byte) (sig []byte, err error) {
	hashed := sha256.Sum256(message)
	sig, err = rsa.SignPKCS1v15(nil, key, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}
	return
}

// Verify that the given message was signed by the private key
// corresponding to the public key we have.
func RSAVerify(pub *rsa.PublicKey, message []byte, sig []byte) bool {
	digest := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying signature: %s\n", err)
		return false
	}
	return true
}

// Reads an existing .pem or rsa keyfile and returns a
// reference to it.
func ReadExistingKey(keyFile string) (*rsa.PrivateKey, error) {
	keyfile, err := os.ReadFile(keyFile)
	CheckErrFatal(err)
	key, err := ssh.ParseRawPrivateKey(keyfile)
	if err != nil {
		err = fmt.Errorf("error parsing key file...%w", err)
		return nil, err
	}
	rsaKey := key.(*rsa.PrivateKey)
	return rsaKey, nil
}

func ParsePublicKey(keyString string) *rsa.PublicKey {
	pKeyBlock, _ := pem.Decode([]byte(keyString))
	if pKeyBlock == nil {
		fmt.Println("Error in pem.Decode, keyblock is nil...")
		panic("Oops")
	}
	pubKeyInterface, err_two := x509.ParsePKIXPublicKey(pKeyBlock.Bytes)
	if err_two != nil {
		fmt.Println("Error in x509.ParsePKIXPublicKey...", err_two)
		panic(err_two)
	}
	pubKey := pubKeyInterface.(*rsa.PublicKey)
	return pubKey
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

func encodePublicKey(key *rsa.PrivateKey) []byte {
	pubKey := key.PublicKey
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		panic(err)
	}
	pemData := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubKeyBytes,
		},
	)
	return pemData
}

func WritePublicKey(key *rsa.PrivateKey, fileName string) {
	bytes := encodePublicKey(key)
	err := os.WriteFile(fileName, bytes, 0600)
	CheckErrFatal(err)
}

func DisplayPublicKey(key *rsa.PrivateKey) {
	pem_key := encodePublicKey(key)
	fmt.Println(string(pem_key))
}

func GenerateRandomKey() *rsa.PrivateKey {
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	CheckErrFatal(err)
	return k
}
