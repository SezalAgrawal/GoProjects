package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

const salt = "ABC123" // can be stored in env

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed string
// as a hex string
func Hash(str string) string {
	// Convert password string to byte slice
	var strBytes = []byte(str)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to str
	strBytes = append(strBytes, salt...)

	// Write str bytes to the hasher
	sha512Hasher.Write(strBytes)

	// Get the SHA-512 hashed str
	var hashedStrBytes = sha512Hasher.Sum(nil)

	// Convert the hashed str to a hex string
	var hashedStrHex = hex.EncodeToString(hashedStrBytes)

	return hashedStrHex
}
