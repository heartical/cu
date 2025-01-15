package cryptography

import (
	"crypto/rand"

	"golang.org/x/crypto/curve25519"
)

// GenerateKeyPair генерирует пару ключей X25519: приватный и публичный.
// Возвращает ошибку, если не удалось сгенерировать случайные данные для приватного ключа.
func GenerateKeyPair() (privateKey, publicKey [32]byte, err error) {
	if _, err = rand.Read(privateKey[:]); err != nil {
		return
	}

	publicKeySlice, err := curve25519.X25519(privateKey[:], curve25519.Basepoint)
	if err != nil {
		return
	}
	copy(publicKey[:], publicKeySlice)

	return
}

// ComputeSharedSecret вычисляет общий секрет на основе приватного и публичного ключей.
// Возвращает ошибку, если не удалось вычислить общий секрет.
func ComputeSharedSecret(privateKey, publicKey [32]byte) (sharedSecret [32]byte, err error) {
	sharedSecretSlice, err := curve25519.X25519(privateKey[:], publicKey[:])
	if err != nil {
		return
	}
	copy(sharedSecret[:], sharedSecretSlice)

	return
}
