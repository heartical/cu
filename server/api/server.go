package api

import (
	"cu/common/cryptography"
	"cu/server/api/database"
	"cu/server/api/router"
	"cu/server/api/security"
	"cu/server/config"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

// Init инициализирует конфигурацию приложения.
// Загружает конфигурацию из файла или переменных окружения.
func Init() {
	config.Load()
	log.Println("Configuration loaded successfully")
}

// startServer запускает HTTP-сервер на указанном порту.
// Инициализирует базу данных, ключи сервера и роутер.
func StartServer(port int) {
	log.Println("Initializing BadgerDB...")
	db, err := database.NewBadgerDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	// Инициализация хранилища ключей сервера.
	log.Println("Initializing server keys storage...")
	serverKeysStorage := security.NewServerKeysStorage(db)

	// Получение или генерация ключей сервера.
	log.Println("Retrieving or generating server keys...")
	serverKeys, err := getOrGenerateServerKeys(serverKeysStorage)
	if err != nil {
		log.Fatalf("Failed to retrieve or generate server keys: %v", err)
	}
	log.Printf("Server public key: %s", hex.EncodeToString(serverKeys.PublicKey[:]))

	// Инициализация роутера с ключами сервера и базой данных.
	log.Println("Setting up router...")
	router := router.NewRouter(serverKeys.PrivateKey, serverKeys.PublicKey, db)

	// Запуск HTTP-сервера.
	log.Printf("Starting server on port :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), router.SetupRoutes()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getOrGenerateServerKeys извлекает ключи сервера из хранилища или генерирует новые.
// Если ключи не найдены в хранилище, они генерируются и сохраняются.
func getOrGenerateServerKeys(storage *security.ServerKeysStorage) (*security.ServerKeys, error) {
	log.Println("Attempting to retrieve server keys from storage...")
	serverKeys, err := storage.Get("production")
	if err != nil {
		log.Println("Server keys not found in storage, generating new keys...")
		serverKeys, err = generateServerKeys()
		if err != nil {
			return nil, fmt.Errorf("failed to generate server keys: %v", err)
		}

		log.Println("Saving generated server keys to storage...")
		if err := storage.Set("production", serverKeys); err != nil {
			return nil, fmt.Errorf("failed to save server keys: %v", err)
		}
		log.Println("Server keys saved successfully")
	}

	log.Println("Server keys retrieved or generated successfully")
	return serverKeys, nil
}

// generateServerKeys генерирует новую пару приватного и публичного ключей.
// Возвращает структуру ServerKeys с ключами или ошибку, если генерация не удалась.
func generateServerKeys() (*security.ServerKeys, error) {
	log.Println("Generating new server keys...")
	privateKey, publicKey, err := cryptography.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	log.Println("Server keys generated successfully")
	return &security.ServerKeys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}
