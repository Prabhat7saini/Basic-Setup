package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

func CompareHashAndPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func EncryptValue(value interface{}, secretKey string) (string, error) {
	key := []byte(secretKey)
	if len(key) != 32 {
		return "", errors.New("secret key must be 32 bytes for AES-256")
	}

	plainBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, plainBytes, nil)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// DecryptValue decrypts a base64-encoded ciphertext and unmarshals it into the target interface.
func DecryptValue(encoded string, secretKey string, value interface{}) error {
	key := []byte(secretKey)
	if len(key) != 32 {
		return errors.New("secret key must be 32 bytes for AES-256")
	}

	cipherText, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherText) < nonceSize {
		return errors.New("cipher text too short")
	}

	nonce, cipherData := cipherText[:nonceSize], cipherText[nonceSize:]
	plainBytes, err := aesGCM.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(plainBytes, value)
}
