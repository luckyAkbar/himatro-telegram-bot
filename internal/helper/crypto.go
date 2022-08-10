package helper

import (
	"crypto/aes"
	"errors"

	"github.com/kumparan/go-utils/encryption"
	"github.com/luckyAkbar/himatro-telegram-bot/internal/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Cryptor :nodoc:
func Cryptor() *encryption.AESCryptor {
	prvKey := config.PrivateKey()
	ivKey := config.IvKey()

	return encryption.NewAESCryptor(prvKey, ivKey, aes.BlockSize)
}

// HashString encrypt given text
func HashString(text string) (string, error) {
	bt, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}

// IsHashedStringMatch check the plain against the cipher using bcrypt.
// If they don't match, will return false
func IsHashedStringMatch(plain, cipher []byte) bool {
	err := bcrypt.CompareHashAndPassword(cipher, plain)
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false
	}
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}
