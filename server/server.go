package main

import (
	"cu/common/assets"
	"cu/common/cryptography"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// //go:embed templates/index.html
// var templateFS embed.FS

// Server представляет серверное приложение.
type Server struct {
	privateKey [32]byte
	publicKey  [32]byte
	sessions   SessionStorage
}

// handleNotFound обрабатывает запросы к несуществующим страницам.
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/game/game-not-found", http.StatusFound)
}

// handleIndex обрабатывает запросы к главной странице.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["room_id"]

	// Используем встроенный файл index.html
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

// handleKeyExchange обрабатывает обмен ключами с клиентом.
func (s *Server) handleKeyExchange(w http.ResponseWriter, r *http.Request) {
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

	baseKey, err := cryptography.ComputeSharedSecret(s.privateKey, clientPublicKey)
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

	s.sessions.SaveSession(&ServerSession{
		AccessKey: accessKey,
		LastUsed:  time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, sessionID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"ServerPublicKey": hex.EncodeToString(s.publicKey[:]),
		"SessionID":       sessionID,
	})
}

// validateEAPI проверяет валидность EAPI для данной сессии.
func (s *Server) validateEAPI(sessionID, receivedEAPI string) bool {
	session, err := s.sessions.GetSession(sessionID)
	if err != nil || time.Now().After(session.ExpiresAt) {
		s.sessions.DeleteSession(sessionID)
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

// handleTunnel обрабатывает запросы через защищенный туннель.
func (s *Server) handleTunnel(callback func(string) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		sessionID := r.FormValue("SessionID")
		eapi := r.FormValue("EAPI")

		if !s.validateEAPI(sessionID, eapi) {
			http.Error(w, "Invalid SessionID or EAPI", http.StatusUnauthorized)
			return
		}

		session, err := s.sessions.GetSession(sessionID)
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
