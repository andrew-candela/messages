package messages

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
)

// Creates a random string and then sign it.
func createSignature(key *rsa.PrivateKey) ([]byte, []byte) {
	token := make([]byte, 64)
	_, _ = rand.Reader.Read(token)
	signed, err := RSASign(key, token)
	if err != nil {
		panic(err)
	}
	return signed, token
}

// Adds the appropriate headers and submit the request.
func makeRequest(r *http.Request, key *rsa.PrivateKey) *http.Response {
	sig, token := createSignature(key)
	sig_str := BytesToString(sig)
	token_str := BytesToString(token)
	r.Header.Set(SIG_HEADER, sig_str)
	r.Header.Set(SIG_TOKEN_HEADER, token_str)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	return res
}

func GetMyIp(host string, key *rsa.PrivateKey) string {
	endpoint := "ip"
	url := fmt.Sprintf("http://%v/%v", host, endpoint)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp := makeRequest(r, key)
	newBuf := make([]byte, 512)
	n, _ := resp.Body.Read(newBuf)
	return string(newBuf[:n])
}

func GetConfig(host string, key *rsa.PrivateKey) string {
	endpoint := "config"
	url := fmt.Sprintf("http://%v/%v", host, endpoint)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		panic(err)
	}
	resp := makeRequest(r, key)
	newBuf := make([]byte, 4096)
	n, _ := resp.Body.Read(newBuf)
	return string(newBuf[:n])
}
