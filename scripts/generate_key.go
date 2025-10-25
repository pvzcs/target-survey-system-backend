package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func main() {
	// Generate 32 bytes for AES-256
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating key: %v\n", err)
		os.Exit(1)
	}

	// Encode to base64 for easy storage
	encoded := base64.StdEncoding.EncodeToString(key)

	fmt.Println("Generated 32-byte encryption key:")
	fmt.Println(encoded)
	fmt.Println("\nAdd this to your .env file:")
	fmt.Printf("ENCRYPTION_KEY=%s\n", encoded)
}
