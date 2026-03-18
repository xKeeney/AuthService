package main

import (
	"auth_service/cmd"
	"flag"
	"log"
)

func main() {
	mode := flag.String("mode", "", "Application mode")

	flag.Parse()

	switch *mode {
	case "server":
		cmd.StartServer()
	case "generate_rsa":
		cmd.GenerateRSAKeys()
	default:
		log.Println("Unexpected mode, run with '-mode=' key.\nAvailable modes:\n  server - start server.\n  generate_rsa - generate private and public keys pair in pem format.")
	}
}
