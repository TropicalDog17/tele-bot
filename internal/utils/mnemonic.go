package utils

import (
	"math/rand"
	"slices"
	"strings"
)

// Store the mnemonic encrypted by a password

func GetEncryptedMnemonic(mnemonic string, password string) (encryptedMnemonic string, salt []byte, err error) {
	// Derive key from password
	key, salt, err := DeriveKey(password)
	if err != nil {
		return "", nil, err
	}
	// Encrypt mnemonic
	encryptedMnemonic, err = Encrypt(key, mnemonic)
	if err != nil {
		return "", nil, err
	}
	return encryptedMnemonic, salt, nil
}

func GetDecryptedMnemonic(key []byte, encryptedMnemonic string) (decryptedMnemonic string, err error) {
	// Decrypt mnemonic
	decryptedMnemonic, err = Decrypt(key, encryptedMnemonic)
	if err != nil {
		return "", err
	}
	return decryptedMnemonic, nil
}

// get 3 random indexes for challenge, ascending order
func GetRandomIndexesForTesting(mnemonicLen int) (indexes [3]int) {
	// Generate 3 random indexes
	for i := 0; i < 3; i++ {
		indexes[i] = rand.Intn(mnemonicLen)
	}
	slices.Sort(indexes[:])
	return indexes
}

func SplitMnemonic(mnemonic string) []string {
	return strings.Split(mnemonic, " ")
}

// from a mnemonic and 3 index of user, check if the provided words are correct
func MnemonicChallenge(mnemonic string, indexes [3]int, providedWords [3]string) (bool, error) {
	// Split mnemonic into words
	words := SplitMnemonic(mnemonic)
	for i, index := range indexes {
		// Check if the word at the index is correct
		if words[index] != providedWords[i] {
			return false, nil
		}
	}
	return true, nil
}

// Return the mnemonic and hole on the indexes
func GenerateMissedWordsMnemonicFromIndexes(mnemonic string, indexes [3]int) string {
	words := SplitMnemonic(mnemonic)
	words[indexes[0]] = "______"
	words[indexes[1]] = "______"
	words[indexes[2]] = "______"
	return strings.Join(words, " ")
}
