package security

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// ServerSession представляет сессию сервера.
type ServerSession struct {
	AccessKey []byte
	LastUsed  time.Time
	ExpiresAt time.Time
}

// SessionStorage предоставляет методы для хранения и извлечения сессий.
type SessionStorage struct {
	db *badger.DB
}

// NewSessionStorage создает новый экземпляр SessionStorage.
func NewSessionStorage(db *badger.DB) *SessionStorage {
	return &SessionStorage{db: db}
}

// SaveSession сохраняет сессию в хранилище.
func (s *SessionStorage) SaveSession(session *ServerSession, id string) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(session); err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry([]byte(id), buf.Bytes()).WithTTL(time.Until(session.ExpiresAt)))
	})
}

// GetSession извлекает сессию по идентификатору id.
func (s *SessionStorage) GetSession(id string) (*ServerSession, error) {
	var session ServerSession
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return gob.NewDecoder(bytes.NewReader(val)).Decode(&session)
		})
	})
	return &session, err
}

// DeleteSession удаляет сессию по идентификатору id.
func (s *SessionStorage) DeleteSession(id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
