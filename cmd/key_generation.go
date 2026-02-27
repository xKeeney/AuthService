package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

const (
	rsaKeyBits = 2048
)

// Generate keys pair for JWT
func GenerateRSAKeys() {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		log.Fatal(err)
	}

	// Private key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privFile.Close()

	if err := pem.Encode(privFile, privatePem); err != nil {
		log.Fatal(err)
	}

	// Public key
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	publicPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	pubFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer pubFile.Close()

	if err := pem.Encode(pubFile, publicPem); err != nil {
		log.Fatal(err)
	}

	log.Println("Keys successfully generated")
}
