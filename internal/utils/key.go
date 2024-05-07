package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

// Derive cipher key from normal password, using PBKDF2
func DeriveKey(password string) (derivedKey []byte, salt []byte, err error) {
	passwordBytes := []byte(password)
	// Generate a random salt
	salt = make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		fmt.Println("Error generating salt:", err)
		return nil, nil, err
	}
	return pbkdf2.Key(passwordBytes, salt, 4096, 32, sha256.New), salt, nil
}

func DeriveKeyFromSalt(password string, salt []byte) ([]byte, error) {
	passwordBytes := []byte(password)
	return pbkdf2.Key(passwordBytes, salt, 4096, 32, sha256.New), nil
}
