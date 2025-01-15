package database

import (
	"github.com/dgraph-io/badger/v4"
)

func NewBadgerDB() (*badger.DB, error) {
	return badger.Open(badger.DefaultOptions("./data"))
}
