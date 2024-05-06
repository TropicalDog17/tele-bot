package utils

import "testing"

func TestEncryptAndDecrypt(t *testing.T) {
	key := []byte("1234567890123456")
	message := "Hello, World!"
	encoded, err := Encrypt(key, message)
	if err != nil {
		t.Errorf("Error encrypting message: %v", err)
	}
	decoded, err := Decrypt(key, encoded)
	if err != nil {
		t.Errorf("Error decrypting message: %v", err)
	}
	if message != decoded {
		t.Errorf("Expected %s, got %s", message, decoded)
	}
}
