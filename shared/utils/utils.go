package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"unicode"

	// "gitlab.com/truemeds-dev-team/truemeds-dev-doctor/truemeds-dev-service/doctorportal-auth-service/shared/constants"
	// "gitlab.com/truemeds-dev-team/truemeds-dev-doctor/truemeds-dev-service/doctorportal-auth-service/shared/constants/exception"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants"
	"github.com/Prabhat7saini/Basic-Setup/shared/constants/exception"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ServiceError[T any](code exception.ErrorCode) constants.ServiceOutput[T] {
	ex := exception.GetException(code)
	return constants.ServiceOutput[T]{
		Exception:      ex,
		Message:        ex.Message,
		HttpStatusCode: ex.HttpStatusCode,
		RespStatusCode: ex.HttpStatusCode,
	}
}
func HandleException[T any](exception constants.Exception) constants.ServiceOutput[T] {
	// ex := exception.GetException(code)
	return constants.ServiceOutput[T]{
		Exception:      &exception,
		Message:        exception.Message,
		HttpStatusCode: exception.HttpStatusCode,
		RespStatusCode: exception.HttpStatusCode,
	}
}

func CheckErrAndLog[T any](
	err error,
	code exception.ErrorCode,
	log *zap.Logger,
	msg string,
) (constants.ServiceOutput[T], bool) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error(msg, zap.Error(err))
		return ServiceError[T](code), true
	}
	return constants.ServiceOutput[T]{}, false
}

func CheckServiceResp(raw []byte, logger *zap.Logger, context string) error {
	var rsp constants.HttpServiceResponse
	if err := json.Unmarshal(raw, &rsp); err != nil {
		logger.Error("cannot unmarshal service response",
			zap.String("context", context),
			zap.Error(err),
			zap.ByteString("body", raw),
		)
		return fmt.Errorf("invalid service response: %w", err)
	}

	logger.Info("service response",
		zap.String("context", context),
		zap.Any("payload", rsp),
	)

	if rsp.Code != http.StatusOK {
		return fmt.Errorf("service returned non‑200 code %d: %s", rsp.Code, rsp.Message)
	}
	return nil
}

func GenerateStrongPassword() (string, error) {
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	digitChars := "0123456789"
	specialChars := "!@#$%^&*()-_=+[]{}|;:,.<>/?"
	allChars := upperChars + lowerChars + digitChars + specialChars

	// Random password length between 8 and 12
	lengthRand, err := rand.Int(rand.Reader, big.NewInt(5)) // 0 to 4
	if err != nil {
		return "", err
	}
	length := 8 + int(lengthRand.Int64()) // 8 to 12

	// Start with one char from each required set
	password := []byte{
		upperChars[randomIndex(len(upperChars))],
		lowerChars[randomIndex(len(lowerChars))],
		digitChars[randomIndex(len(digitChars))],
		specialChars[randomIndex(len(specialChars))],
	}

	// Fill the rest with random characters from all sets
	for len(password) < length {
		password = append(password, allChars[randomIndex(len(allChars))])
	}

	// Shuffle the password
	for i := range password {
		j := randomIndex(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	return string(password), nil
}

func randomIndex(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func UsernameToEmail(username *string) *string {
	if username == nil || *username == "" {
		return nil
	}
	email := *username + constants.TruemedsEmailExtension
	return &email
}

func CapitalizeWords(s string) string {
	if s == "" {
		return ""
	}
	words := strings.Split(s, " ")
	for i, word := range words {
		if word == "" {
			continue
		}
		runes := []rune(strings.ToLower(word))
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}
