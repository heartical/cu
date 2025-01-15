package controllers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"cu/common/assets"
	"cu/common/cryptography"
	"cu/server/api/security"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// PlayController представляет контроллер для обработки запросов.
type PlayController struct {
	privateKey [32]byte
	publicKey  [32]byte
	sessions   *security.SessionStorage
	db         *badger.DB
}

// NewPlayController создает новый контроллер.
func NewPlayController(privateKey, publicKey [32]byte, db *badger.DB) *PlayController {
	return &PlayController{
		privateKey: privateKey,
		publicKey:  publicKey,
		sessions:   security.NewSessionStorage(db),
		db:         db,
	}
}

// PageRequest обрабатывает запрос на отображение страницы.
func (pc *PlayController) PageRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	tmpl, err := template.ParseFS(assets.TemplateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, roomID); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// KeyExchangeRequest обрабатывает обмен ключами с клиентом.
func (pc *PlayController) KeyExchangeRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	clientPublicKeyHex := r.FormValue("ClientPublicKey")
	clientPublicKeyBytes, err := hex.DecodeString(clientPublicKeyHex)
	if err != nil {
		http.Error(w, "Unable to decode ClientPublicKey", http.StatusBadRequest)
		return
	}

	var clientPublicKey [32]byte
	copy(clientPublicKey[:], clientPublicKeyBytes)

	baseKey, err := cryptography.ComputeSharedSecret(pc.privateKey, clientPublicKey)
	if err != nil {
		http.Error(w, "Unable to compute shared secret", http.StatusInternalServerError)
		return
	}

	sessionID := uuid.New().String()
	sessionKey, err := cryptography.GenerateSessionKey(baseKey[:], sessionID)
	if err != nil {
		http.Error(w, "Unable to generate session key", http.StatusInternalServerError)
		return
	}

	accessKey, err := cryptography.GenerateAccessKey(sessionKey)
	if err != nil {
		http.Error(w, "Unable to generate access key", http.StatusInternalServerError)
		return
	}

	pc.sessions.SaveSession(&security.ServerSession{
		AccessKey: accessKey,
		LastUsed:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, sessionID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"ServerPublicKey": hex.EncodeToString(pc.publicKey[:]),
		"SessionID":       sessionID,
	})
}

// validateEAPI проверяет валидность EAPI для данной сессии.
func (pc *PlayController) validateEAPI(sessionID, receivedEAPI string) bool {
	session, err := pc.sessions.GetSession(sessionID)
	if err != nil || time.Now().After(session.ExpiresAt) {
		pc.sessions.DeleteSession(sessionID)
		return false
	}

	currentTimestamp := time.Now().Unix() / 30
	expectedEAPI := hex.EncodeToString(cryptography.ComputeEAPI(session.AccessKey, currentTimestamp))
	previousEAPI := hex.EncodeToString(cryptography.ComputeEAPI(session.AccessKey, currentTimestamp-1))

	if receivedEAPI != expectedEAPI && receivedEAPI != previousEAPI {
		return false
	}

	session.LastUsed = time.Now()
	return true
}

// ActionRequest обрабатывает запросы через защищенный туннель.
func (pc *PlayController) ActionRequest(callback func(string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		sessionID := r.FormValue("SessionID")
		eapi := r.FormValue("EAPI")

		if !pc.validateEAPI(sessionID, eapi) {
			http.Error(w, "Invalid SessionID or EAPI", http.StatusUnauthorized)
			return
		}

		session, err := pc.sessions.GetSession(sessionID)
		if err != nil {
			http.Error(w, "Invalid SessionID", http.StatusUnauthorized)
			return
		}

		data := r.FormValue("Data")

		decrypted, err := cryptography.DecryptAES(data, session.AccessKey)
		if err != nil {
			http.Error(w, "Unable to decrypt data", http.StatusBadRequest)
			return
		}

		decryptedString := string(decrypted)

		response := callback(decryptedString)

		encrypted, err := cryptography.EncryptAES([]byte(response), session.AccessKey)
		if err != nil {
			http.Error(w, "Unable to encrypt response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"Data": string(encrypted)})
	}
}
