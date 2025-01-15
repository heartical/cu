package e2e

import (
	"cu/common/cryptography"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"syscall/js"
	"time"
)

// Client представляет клиента для взаимодействия с сервером.
type Client struct {
	PrivateKey [32]byte
	PublicKey  [32]byte
	SharedKey  [32]byte
	DerivedKey []byte
	AccessKey  []byte
	SessionID  string
	ServerURL  string
}

// NewClient создает новый клиент с указанным URL сервера.
func NewClient(serverURL string) (*Client, error) {
	privateKey, publicKey, err := cryptography.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate client keys: %w", err)
	}

	return &Client{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		ServerURL:  serverURL,
	}, nil
}

// LoadSessionFromLocalStorage загружает сессию из локального хранилища.
func (c *Client) LoadSessionFromLocalStorage() bool {
	localStorage := js.Global().Get("localStorage")
	sessionIDVal := localStorage.Call("getItem", "SessionID")
	if sessionIDVal.IsNull() || sessionIDVal.IsUndefined() {
		return false
	}
	accessKeyHexVal := localStorage.Call("getItem", "AccessKey")
	if accessKeyHexVal.IsNull() || accessKeyHexVal.IsUndefined() {
		return false
	}

	accessKey, err := hex.DecodeString(accessKeyHexVal.String())
	if err != nil {
		log.Printf("failed to decode AccessKey: %v", err)
		return false
	}

	c.SessionID = sessionIDVal.String()
	c.AccessKey = accessKey

	return true
}

// saveSessionToLocalStorage сохраняет сессию в локальное хранилище.
func (c *Client) saveSessionToLocalStorage() {
	localStorage := js.Global().Get("localStorage")
	localStorage.Call("setItem", "SessionID", c.SessionID)
	localStorage.Call("setItem", "AccessKey", hex.EncodeToString(c.AccessKey))
}

// ExchangeKeysWithServer выполняет обмен ключами с сервером.
func (c *Client) ExchangeKeysWithServer() error {
	clientPublicKeyHex := hex.EncodeToString(c.PublicKey[:])
	resp, err := http.PostForm(c.ServerURL+"/key-exchange", url.Values{
		"ClientPublicKey": {clientPublicKeyHex},
	})
	if err != nil {
		return fmt.Errorf("failed to send public key: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		ServerPublicKey string `json:"ServerPublicKey"`
		SessionID       string `json:"SessionID"`
		Error           string `json:"Error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse server response: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf("server returned an error: %s", result.Error)
	}

	serverPublicKey, err := hex.DecodeString(result.ServerPublicKey)
	if err != nil {
		return fmt.Errorf("failed to decode server public key: %w", err)
	}

	c.SharedKey, err = cryptography.ComputeSharedSecret(c.PrivateKey, *(*[32]byte)(serverPublicKey))
	if err != nil {
		return fmt.Errorf("failed to create SharedKey: %w", err)
	}

	c.SessionID = result.SessionID

	c.DerivedKey, err = cryptography.GenerateSessionKey(c.SharedKey[:], c.SessionID)
	if err != nil {
		return fmt.Errorf("failed to derive DerivedKey: %w", err)
	}

	c.AccessKey, err = cryptography.GenerateAccessKey(c.DerivedKey)
	if err != nil {
		return fmt.Errorf("failed to derive AccessKey: %w", err)
	}

	c.saveSessionToLocalStorage()
	return nil
}

// GetCurrentEAPI возвращает текущий EAPI (Encrypted API Key).
func (c *Client) GetCurrentEAPI() string {
	timestamp := time.Now().Unix() / 30
	eapi := cryptography.ComputeEAPI(c.AccessKey, timestamp)
	return hex.EncodeToString(eapi)
}

// SendMessageToServer отправляет зашифрованное сообщение на сервер.
func (c *Client) SendMessageToServer(message string) (string, error) {
	encrypted, err := cryptography.EncryptAES([]byte(message), c.AccessKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	resp, err := http.PostForm(c.ServerURL+"/action", url.Values{
		"Data":      {encrypted},
		"EAPI":      {c.GetCurrentEAPI()},
		"SessionID": {c.SessionID},
	})
	if err != nil {
		return "", fmt.Errorf("failed to send encrypted message: %w", err)
	}
	defer resp.Body.Close()

	var response struct {
		Data  string `json:"Data"`
		Error string `json:"Error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to parse server response: %w", err)
	}

	if response.Error != "" {
		return "", fmt.Errorf("server returned an error: %s", response.Error)
	}

	decrypted, err := cryptography.DecryptAES(response.Data, c.AccessKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt server response: %w", err)
	}

	return decrypted, nil
}
