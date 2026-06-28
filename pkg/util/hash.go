package util

import (
	"crypto/md5"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const HashingCost int = 12 // if under 4, it will return error, see golang.org/x/crypto/bcrypt/bcrypt.go:289

type (
	CryptoHashService interface {
		CreatePasswordHash(plainPassword string) (hashedPassword string, err error)
		ValidatePassword(hashedPassword, plainPassword string) (isValid bool)
		CreateMD5Hash(plainText string) (hashedText string)
	}
	cryptoHashService struct {
		bcryptHashingCost int
	}
)

func NewCryptoHash() CryptoHashService {
	return &cryptoHashService{
		bcryptHashingCost: HashingCost,
	}
}

// CreatePasswordHash creates a password hash of given `plainPassword`
func (s *cryptoHashService) CreatePasswordHash(plainPassword string) (hashedPassword string, err error) {
	passwordHashInBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), s.bcryptHashingCost)
	if err != nil {
		return
	}
	hashedPassword = string(passwordHashInBytes)
	return
}

// ValidatePassword validates given `hashedPassword` against `plainPassword`. It returns true if given passwords are matched.
func (s *cryptoHashService) ValidatePassword(hashedPassword, plainPassword string) (isValid bool) {
	hashedPasswordInBytes := []byte(hashedPassword)
	plainPasswordInBytes := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordInBytes, plainPasswordInBytes)
	isValid = err == nil
	return
}

// CreateMD5Hash returns md5 hash value of `plainText`
func (s *cryptoHashService) CreateMD5Hash(plainText string) (hashedText string) {
	strInByte := []byte(plainText)
	resultInByte := md5.Sum(strInByte)
	hashedText = fmt.Sprintf("%x", resultInByte)
	return
}
