package messages

import (
	"crypto"
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

func RSASign(key *rsa.PrivateKey, message []byte) (sig []byte, err error) {
	hashed := sha256.Sum256(message)
	sig, err = rsa.SignPKCS1v15(nil, key, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}
	return
}

func RSAVerify(pub *rsa.PublicKey, message []byte, sig []byte) bool {
	digest := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, digest[:], sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying signature: %s\n", err)
		return false
	}
	return true
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

// https://stackoverflow.com/a/70719783/14223687 for a good example
func ParsePublicKey(keyString string) rsa.PublicKey {

	var spkiPem = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoZ67dtUTLxoXnNEzRBFB
mwukEJGC+y69cGgpNbtElQj3m4Aft/7cu9qYbTNguTSnCDt7uovZNb21u1vpZwKH
yVgFEGO4SA8RNnjhJt2D7z8RDMWX3saody7jo9TKlrPABLZGo2o8vadW8Dly/v+I
d0YDheCkVCoCEeUjQ8koXZhTwhYkGPu+vkdiqX5cUaiVTu1uzt591aO5Vw/hV4DI
hFKnOTnYXnpXiwRwtPyYoGTa64yWfi2t0bv99qz0BgDjQjD0civCe8LRXGGhyB1U
1aHjDDGEnulTYJyEqCzNGwBpzEHUjqIOXElFjt55AFGpCHAuyuoXoP3gQvoSj6RC
sQIDAQAB
-----END PUBLIC KEY-----`
	pKeyBlock, _ := pem.Decode([]byte(spkiPem))
	if pKeyBlock == nil {
		fmt.Println("Error in pem.Decode...")
		panic("Oops")
	}
	pubKeyInterface, err_two := x509.ParsePKIXPublicKey(pKeyBlock.Bytes)
	if err_two != nil {
		fmt.Println("Error in x509.ParsePKIXPublicKey...", err_two)
		panic(err_two)
	}
	pubKey := pubKeyInterface.(*rsa.PublicKey)
	return *pubKey
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
