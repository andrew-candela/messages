package main

import (
	"crypto/rsa"
	"flag"
	"os"

	"github.com/andrew-candela/messages/messages"
)

// this is just for testing
func getKey() *rsa.PublicKey {
	k, _ := messages.ReadExistingKey("/Users/andrewcandela/.ssh/id_rsa")
	return &k.PublicKey
}

func makeTargets(k *rsa.PublicKey) *[]messages.RecipientDetails {
	return &[]messages.RecipientDetails{
		{
			DestinationHostPort: "10.0.0.176:1053",
			PublicKey:           *k,
		},
		{
			DestinationHostPort: "aim.andrewcandela.com:1053",
			PublicKey:           *k,
		},
	}
}

func main() {
	var host string
	var user string
	flag.StringVar(&host, "host", "10.0.0.186:1053", "Provide the HOST:PORT to write to")
	flag.StringVar(&user, "user", os.Getenv("USER"), "What do you call yourself?")
	flag.Parse()
	k := getKey()
	targets := makeTargets(k)
	messages.ProduceMessages(*targets, user)
}
