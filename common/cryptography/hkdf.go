package cryptography

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"io"

	"golang.org/x/crypto/hkdf"
)

// deriveKey генерирует ключ заданной длины на основе секрета, соли и дополнительной информации.
// Используется алгоритм HKDF для безопасного расширения ключа.
func deriveKey(secret, salt, info []byte, keyLen int) ([]byte, error) {
	hkdfReader := hkdf.New(sha256.New, secret, salt, info)
	key := make([]byte, keyLen)

	if _, err := io.ReadFull(hkdfReader, key); err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateSessionKey генерирует сессионный ключ на основе базового ключа и API-ключа.
func GenerateSessionKey(baseKey []byte, apiKey string) ([]byte, error) {
	return deriveKey(baseKey, []byte(apiKey), []byte("SessionKey"), 32)
}

// GenerateAccessKey генерирует ключ доступа на основе сессионного ключа.
func GenerateAccessKey(sessionKey []byte) ([]byte, error) {
	return deriveKey(sessionKey, nil, []byte("AccessKey"), 32)
}

// ComputeEAPI вычисляет HMAC-SHA256 на основе ключа доступа и временной метки.
func ComputeEAPI(accessKey []byte, timestamp int64) []byte {
	h := hmac.New(sha256.New, accessKey)
	binary.Write(h, binary.BigEndian, timestamp)
	return h.Sum(nil)
}
