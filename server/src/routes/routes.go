package routing

import (
	mux "github.com/gorilla/mux"

	controllers "akshayGudhate/whatsapp-bridge/src/controllers"
)

func Routes() *mux.Router {
	// router
	router := mux.NewRouter()

	// invoice
	router.HandleFunc("/api/whatsapp/send", controllers.SendMessage).Methods("POST")
	router.HandleFunc("/api/connect/qr", controllers.GetConnectionQRCode).Methods("GET")

	return router
}
