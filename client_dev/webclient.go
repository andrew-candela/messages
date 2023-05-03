package main

import (
	"fmt"
	"log"

	"github.com/andrew-candela/messages/messages"
)

func main() {
	// start webclient
	fmt.Println("Starting webclient")
	log.SetFlags(0)
	messages.ExampleClient()
}
