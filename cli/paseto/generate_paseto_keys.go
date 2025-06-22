package main

import (
	"fmt"

	paseto "aidanwoods.dev/go-paseto"
)

func main() {
	// Generate new V4 secret key
	secretKey := paseto.NewV4AsymmetricSecretKey()

	// Extract public key from the secret key
	publicKey := secretKey.Public()

	// Export keys as hex strings
	fmt.Println("Private Key (keep this safe):")
	fmt.Println(secretKey.ExportHex())

	fmt.Println("\nPublic Key (can be shared):")
	fmt.Println(publicKey.ExportHex())
}
