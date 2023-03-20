package main

import (
	"flag"
	"os"

	"github.com/andrew-candela/messages/messages"
)

const (
	HOST = "10.0.0.186"
	PORT = "1053"
	TYPE = "udp"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "10.0.0.186:1053", "Provide the HOST:PORT to write to")
	flag.Parse()
	messages.ProduceMessages(host, os.Getenv("USER"))
}
