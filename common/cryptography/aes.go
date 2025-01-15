package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// EncryptAES шифрует данные с использованием AES-GCM.
// Возвращает зашифрованный текст в виде hex-строки.
func EncryptAES(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %w", err)
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return hex.EncodeToString(cipherText), nil
}

// DecryptAES расшифровывает данные, зашифрованные с использованием AES-GCM.
// Принимает зашифрованный текст в виде hex-строки и возвращает расшифрованный текст.
func DecryptAES(hexCipherText string, key []byte) (string, error) {
	cipherText, err := hex.DecodeString(hexCipherText)
	if err != nil {
		return "", fmt.Errorf("error decoding ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("error creating AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting message: %w", err)
	}

	return string(plainText), nil
}
