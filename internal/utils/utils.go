package utils

import "crypto/ecdsa"

type UtilsInterface interface {
	GetEncryptedMnemonic(mnemonic, password string) (string, []byte, error)
	MnemonicChallenge(mnemonic string, indexes [3]int, providedWords [3]string) (bool, error)
	SplitMnemonic(mnemonic string) []string
	GenerateMnemonic() (string, error)
}

type Utils struct{}

func (u Utils) GetEncryptedMnemonic(mnemonic, password string) (string, []byte, error) {
	return GetEncryptedMnemonic(mnemonic, password)
}

func (u Utils) MnemonicChallenge(mnemonic string, indexes [3]int, providedWords [3]string) (bool, error) {
	return MnemonicChallenge(mnemonic, indexes, providedWords)
}

func (u Utils) SplitMnemonic(mnemonic string) []string {
	return SplitMnemonic(mnemonic)
}

func (u Utils) DerivePrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	return DerivePrivateKeyFromMnemonic(mnemonic)
}

func (u Utils) GenerateMnemonic() (string, error) {
	return Generate24wordsRandomMnemonic()
}
