package router

import (
	"cu/server/api/controllers"
	"net/http"

	"github.com/dgraph-io/badger/v4"
	"github.com/gorilla/mux"
)

type Router struct {
	muxRouter  *mux.Router
	privateKey [32]byte
	publicKey  [32]byte
	db         *badger.DB
}

func NewRouter(privateKey, publicKey [32]byte, db *badger.DB) *Router {
	return &Router{
		muxRouter:  mux.NewRouter().StrictSlash(true),
		privateKey: privateKey,
		publicKey:  publicKey,
		db:         db,
	}
}

func (router *Router) SetupRoutes() *mux.Router {
	// Инициализация контроллеров с передачей ключей и базы данных
	playController := controllers.NewPlayController(router.privateKey, router.publicKey, router.db)

	router.muxRouter.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))),
	)

	// Настройка маршрутов
	router.muxRouter.HandleFunc("/{room_id}", playController.PageRequest).Methods("GET")
	router.muxRouter.HandleFunc("/key-exchange", playController.KeyExchangeRequest).Methods("POST")
	router.muxRouter.HandleFunc("/action", playController.ActionRequest(func(message string) string {
		if message == "ping" {
			return "pong"
		}
		return message
	})).Methods("POST")

	return router.muxRouter
}
