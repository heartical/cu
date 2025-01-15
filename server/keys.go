package main

import (
	"encoding/hex"
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

// ServerKeys представляет приватный и публичный ключи сервера.
type ServerKeys struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
}

// ServerKeysStorage предоставляет методы для хранения и извлечения ключей сервера.
type ServerKeysStorage struct {
	db *badger.DB
}

// NewServerKeysStorage создает новый экземпляр ServerKeysStorage.
func NewServerKeysStorage(db *badger.DB) *ServerKeysStorage {
	return &ServerKeysStorage{db: db}
}

// Get извлекает ключи сервера по идентификатору keyID.
func (s *ServerKeysStorage) Get(keyID string) (*ServerKeys, error) {
	var keys ServerKeys
	err := s.db.View(func(txn *badger.Txn) error {
		privItem, err := txn.Get([]byte(fmt.Sprintf("%s:privateKey", keyID)))
		if err != nil {
			return err
		}
		pubItem, err := txn.Get([]byte(fmt.Sprintf("%s:publicKey", keyID)))
		if err != nil {
			return err
		}
		privValue, err := privItem.ValueCopy(nil)
		if err != nil {
			return err
		}
		privBytes, err := hex.DecodeString(string(privValue))
		if err != nil {
			return err
		}
		pubValue, err := pubItem.ValueCopy(nil)
		if err != nil {
			return err
		}
		pubBytes, err := hex.DecodeString(string(pubValue))
		if err != nil {
			return err
		}
		copy(keys.PrivateKey[:], privBytes)
		copy(keys.PublicKey[:], pubBytes)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &keys, nil
}

// Set сохраняет ключи сервера по идентификатору keyID.
func (s *ServerKeysStorage) Set(keyID string, keys *ServerKeys) error {
	return s.db.Update(func(txn *badger.Txn) error {
		privHex := hex.EncodeToString(keys.PrivateKey[:])
		pubHex := hex.EncodeToString(keys.PublicKey[:])
		if err := txn.Set([]byte(keyID+":privateKey"), []byte(privHex)); err != nil {
			return err
		}
		return txn.Set([]byte(keyID+":publicKey"), []byte(pubHex))
	})
}

// Delete удаляет ключи сервера по идентификатору keyID.
func (s *ServerKeysStorage) Delete(keyID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(keyID + ":privateKey"))
	})
}
