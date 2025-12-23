package main

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func main() {
	key := paseto.NewV4AsymmetricSecretKey()
	fmt.Printf("Hex Key: %s\n", key.ExportHex())
}
