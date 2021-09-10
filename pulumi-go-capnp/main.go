package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("[ERROR] required subcommand: 'server' or 'client'")
	}
	switch os.Args[1] {
	case "server":
		serverMain()
	case "client":
		clientMain()
	default:
		log.Fatalf("[ERROR] required subcommand: 'server' or 'client'")
	}
}
