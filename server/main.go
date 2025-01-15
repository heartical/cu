package main

import (
	"cu/common/cryptography"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	db, err := badger.Open(badger.DefaultOptions("./data"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	serverKeysStorage := NewServerKeysStorage(db)
	serverKeys, err := serverKeysStorage.Get("production")
	if err != nil {
		serverKeys, err = generateServerKeys()
		if err != nil {
			panic(err)
		}
		if err := serverKeysStorage.Set("production", serverKeys); err != nil {
			panic(err)
		}
	}

	log.Printf("Public Key: %s", hex.EncodeToString(serverKeys.PublicKey[:]))

	srv := Server{
		privateKey: serverKeys.PrivateKey,
		publicKey:  serverKeys.PublicKey,
		sessions:   *NewSessionStorage(db),
	}

	r.HandleFunc("/exchange-keys", srv.handleKeyExchange).Methods("POST")
	r.HandleFunc("/action", srv.handleTunnel(func(data string) string {
		if data == "ping" {
			return "pong"
		}
		return data
	})).Methods("POST")

	r.HandleFunc("/", srv.handleNotFound).Methods("GET")
	r.HandleFunc("/game/{room_id}", srv.handleIndex).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(srv.handleNotFound)

	log.Println("Starting server on :8010")
	if err := http.ListenAndServe(":8010", r); err != nil {
		panic(err)
	}
}

// generateServerKeys генерирует новые ключи сервера.
func generateServerKeys() (*ServerKeys, error) {
	privateKey, publicKey, err := cryptography.GenerateKeyPair()
	if err != nil {
		return nil, err
	}
	return &ServerKeys{PrivateKey: privateKey, PublicKey: publicKey}, nil
}
