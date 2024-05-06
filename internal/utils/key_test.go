package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeriveKey(t *testing.T) {
	password := "password"
	key, _, err := DeriveKey(password)
	if err != nil {
		t.Errorf("Error deriving key: %v", err)
	}
	if key == nil {
		t.Error("Key is nil")
	}
	fmt.Println(key)
	assert.Len(t, key, 32)
}
