package api

import (
	http "net/http"

	handlers "github.com/gorilla/handlers"
	mux "github.com/gorilla/mux"
)

func GetHandlerWithRoutes() http.Handler {
	// router
	router := mux.NewRouter()

	// invoice
	router.HandleFunc("/api/whatsapp/send", SendMessage).Methods("POST")
	router.HandleFunc("/api/connect/qr", GetConnectionQRCode).Methods("GET")

	// cors handler
	handlerWithCors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)(router)

	return handlerWithCors
}
