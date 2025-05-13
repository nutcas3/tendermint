package auth

import (
	"crypto/sha256"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateMessageSignature verifies the digital signature of a message
func ValidateMessageSignature(message []byte, signature []byte, pubKey sdk.AccAddress) error {
	// Implement digital signature verification logic
	// to do: use the public key to verify the signature
	// This is a placeholder for the actual implementation
	// You would typically use a cryptographic library to verify the signature
	return nil
}

// HashMessage returns the SHA-256 hash of a message
func HashMessage(message []byte) []byte {
	hash := sha256.Sum256(message)
	return hash[:]
}

// ValidateMessageIntegrity checks if the message hash matches the expected hash
func ValidateMessageIntegrity(message []byte, expectedHash []byte) error {
	actualHash := HashMessage(message)
	if !equal(actualHash, expectedHash) {
		return errors.New("message integrity check failed")
	}
	return nil
}

// equal checks if two byte slices are equal
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
