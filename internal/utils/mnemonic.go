package utils

import (
	"math/rand"
	"slices"
	"strings"

	"github.com/cosmos/go-bip39"
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

// get 3 random indexes for challenge, ascending order, unique
func GetRandomIndexesForTesting(mnemonicLen int) (indexes [3]int) {
	// Get 3 random indexes
	indexes[0] = rand.Intn(mnemonicLen)
	indexes[1] = rand.Intn(mnemonicLen)
	indexes[2] = rand.Intn(mnemonicLen)
	// Sort indexes
	slices.Sort(indexes[:])
	// Check if indexes are unique
	for i := 0; i < 2; i++ {
		if indexes[i] == indexes[i+1] {
			// If not unique, get new indexes
			return GetRandomIndexesForTesting(mnemonicLen)
		}
	}
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

func Generate24wordsRandomMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
