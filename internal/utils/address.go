package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

const (
	defaultDerivationPath = "m/44'/60'/0'/0/0"
	defaultBech32Prefix   = "inj"
)

func DerivePublicKeyFromMnemonic(mnemonic string) (string, error) {
	var wallet *hdwallet.Wallet
	var err error

	wallet, err = hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	path := hdwallet.MustParseDerivationPath(defaultDerivationPath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}

	ethereumAddress := account.Address.Hex()
	fmt.Println("Ethereum Address:", ethereumAddress)

	// Decode hex string to bytes, skipping the '0x' prefix
	addressBytes, err := hex.DecodeString(ethereumAddress[2:])
	if err != nil {
		return "", err
	}

	// Convert to Bech32 without hashing
	injectiveAddress, err := bech32.ConvertAndEncode(defaultBech32Prefix, addressBytes)
	if err != nil {
		return "", err
	}

	fmt.Println("Injective Address:", injectiveAddress)
	return injectiveAddress, nil
}

type PrivateKey struct {
	PrivateKey *ecdsa.PrivateKey
}

func (pk *PrivateKey) Type() string {
	return "ethsecp256k1"
}

func DerivePrivateKeyFromMnemonic(mnemonic string) (*ecdsa.PrivateKey, error) {
	var wallet *hdwallet.Wallet
	var err error

	wallet, err = hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path := hdwallet.MustParseDerivationPath(defaultDerivationPath)
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func ExportEncryptedPrivateKeyFromMnemonicAndPassphrase(mnemonic, passphrase string) (string, error) {
	privateKey, err := DerivePrivateKeyFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}

	// Encrypt private key
	encryptedPrivateKey, err := Encrypt([]byte(passphrase), ECDSAToString(privateKey))
	if err != nil {
		return "", err
	}

	return encryptedPrivateKey, nil
}

func DerivePrivateKeyFromEncryptedAndPassphrase(encryptedPrivateKey, passphrase string) (*ecdsa.PrivateKey, error) {
	// Decrypt private key
	decryptedPrivateKey, err := Decrypt([]byte(passphrase), encryptedPrivateKey)
	if err != nil {
		return nil, err
	}

	// Decode hex string to bytes
	privateKeyBytes, err := hex.DecodeString(decryptedPrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := hexToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func hexToECDSA(hexKey []byte) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hex.EncodeToString(hexKey))
}

func ECDSAToString(key *ecdsa.PrivateKey) string {
	return hex.EncodeToString(key.D.Bytes())
}
