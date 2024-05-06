package utils

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// from: pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear
// public address: inj1wfawuv6fslzjlfa4v7exv27mk6rpfeyvhvxchc
var privateKey = "6C212553111B370A8FFDC682954495B7B90A73CEDAB7106323646A4F2C4E668F"
var mnemonic = "pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"
var password = "1234567890123456"
var injectiveAddress = "inj1wfawuv6fslzjlfa4v7exv27mk6rpfeyvhvxchc"

func TestDeriveAddressFromMnemonic(t *testing.T) {

	address, err := DerivePublicKeyFromMnemonic(mnemonic)
	if err != nil {
		t.Errorf("Error deriving address from mnemonic: %v", err)
	}
	if address != injectiveAddress {
		t.Errorf("Expected %s, got %s", injectiveAddress, address)
	}
}

func TestDerivePrivateKeyFromMnemonic(t *testing.T) {
	privateKey := "6C212553111B370A8FFDC682954495B7B90A73CEDAB7106323646A4F2C4E668F"
	mnemonic := "pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"
	derivedPrivateKey, err := DerivePrivateKeyFromMnemonic(mnemonic)
	derivedPrivateKeyString := hex.EncodeToString(crypto.FromECDSA(derivedPrivateKey))

	if err != nil {
		t.Errorf("Error deriving private key from mnemonic: %v", err)
	}
	if strings.ToUpper(derivedPrivateKeyString) != privateKey {
		t.Errorf("Expected %s, got %s", privateKey, derivedPrivateKey)
	}
}

func TestEncryptAndDecryptPrivateKey(t *testing.T) {
	exportedPrivateKey, err := ExportEncryptedPrivateKeyFromMnemonicAndPassphrase(mnemonic, password)
	require.NoError(t, err)

	derivedPrivateKey, err := DerivePrivateKeyFromEncryptedAndPassphrase(exportedPrivateKey, password)
	require.NoError(t, err)

	derivedPrivateKeyString := hex.EncodeToString(crypto.FromECDSA(derivedPrivateKey))
	require.Equal(t, privateKey, strings.ToUpper(derivedPrivateKeyString))
}
