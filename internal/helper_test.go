package internal

import (
	"strings"
	"testing"

	"context"

	"github.com/TropicalDog17/tele-bot/internal/utils"
	"github.com/awnumar/memguard"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetrievePrivateKeyFromRedis_E2E(t *testing.T) {
	// Set up Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Update with your Redis instance address
	})
	defer redisClient.Close()

	// Set up test data
	username := "testuser"
	password := "tropical"
	expectedPrivateKey := "6C212553111B370A8FFDC682954495B7B90A73CEDAB7106323646A4F2C4E668F"
	mnemonic := "pony glide frown crisp unfold lawn cup loan trial govern usual matrix theory wash fresh address pioneer between meadow visa buffalo keep gallery swear"
	salt := "testsalt"

	// Encrypt mnemonic and store in Redis
	key, err := utils.DeriveKeyFromSalt(password, []byte(salt))
	require.NoError(t, err)
	encryptedMnemonic, err := utils.Encrypt(key, mnemonic)
	require.NoError(t, err)
	err = redisClient.HSet(context.Background(), username, "encryptedMnemonic", encryptedMnemonic).Err()
	require.NoError(t, err)
	err = redisClient.HSet(context.Background(), username, "salt", salt).Err()
	require.NoError(t, err)

	lockedPassword := memguard.NewBufferFromBytes([]byte(password))
	// Call the function being tested
	privateKey, err := RetrievePrivateKeyFromRedis(redisClient, username, lockedPassword)
	require.NoError(t, err)

	// Assert that the private key is not nil
	assert.NotNil(t, privateKey)

	assert.Equal(t, expectedPrivateKey, strings.ToUpper(privateKey.String()))
	// Clean up
	privateKey.Destroy()

}
