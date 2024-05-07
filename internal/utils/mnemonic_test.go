package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEncrpytedMnemonicAndDecrypt(t *testing.T) {
	mnemonic := "abandon"
	password := "password"

	// Encrypt mnemonic
	encryptedMnemonic, salt, err := GetEncryptedMnemonic(mnemonic, password)
	if err != nil {
		t.Errorf("Error getting encrypted mnemonic: %v", err)
	}

	// Decrypt mnemonic
	key, err := DeriveKeyFromSalt(password, salt)
	if err != nil {
		t.Errorf("Error deriving key from salt: %v", err)
	}
	decryptedMnemonic, err := GetDecryptedMnemonic(key, encryptedMnemonic)
	if err != nil {
		t.Errorf("Error getting decrypted mnemonic: %v", err)
	}
	require.Equal(t, mnemonic, decryptedMnemonic)
}

func TestGetRandomIndexesForTesting(t *testing.T) {
	indexes := GetRandomIndexesForTesting(len(SplitMnemonic(testMnemonic)))
	require.Len(t, indexes, 3)
	for _, index := range indexes {
		require.True(t, index >= 0 && index < len(SplitMnemonic(testMnemonic)))
	}
}

func TestChallengeMnemonic(t *testing.T) {
	mockIndexes := [3]int{2, 4, 7}
	// frown unfold loan
	mockProvidedWords := [3]string{"frown", "unfold", "loan"}
	result, err := MnemonicChallenge(testMnemonic, mockIndexes, mockProvidedWords)
	require.NoError(t, err)
	require.True(t, result)
}

func TestGenerateMnemonic(t *testing.T) {
	mnemonic, err := Generate24wordsRandomMnemonic()
	require.NoError(t, err)
	require.NotEmpty(t, mnemonic)
	require.Len(t, SplitMnemonic(mnemonic), 24)
}
