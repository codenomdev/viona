package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const PasswordString string = "Password@1234"

func TestNew(t *testing.T) {
	assert.NotNil(t, NewCryptoHash())
}

func TestCreatePasswordHash(t *testing.T) {
	hash := NewCryptoHash()
	hashed, err := hash.CreatePasswordHash(PasswordString)

	assert.NoError(t, err)
	assert.NotNil(t, hashed)
}

func TestCreateMD5Hash(t *testing.T) {
	hash := NewCryptoHash()
	hashed := hash.CreateMD5Hash(PasswordString)
	assert.NotNil(t, hashed)
}

func TestValidatePasswordHash(t *testing.T) {
	hash := NewCryptoHash()
	hashed, err := hash.CreatePasswordHash(PasswordString)
	valid := hash.ValidatePassword(hashed, PasswordString)
	assert.NoError(t, err)
	assert.NotNil(t, hashed)
	assert.True(t, valid)
}

func TestUnvalidatedPasswordHash(t *testing.T) {
	hash := NewCryptoHash()
	hashed, err := hash.CreatePasswordHash("PasswordString")
	valid := hash.ValidatePassword(hashed, PasswordString)
	assert.NoError(t, err)
	assert.NotNil(t, hashed)
	assert.False(t, valid)
}
