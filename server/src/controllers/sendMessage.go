package controllers

import (
	json "encoding/json"
	http "net/http"

	services "akshayGudhate/whatsapp-bridge/src/services"
)

//////////////////
//    struct    //
//////////////////

type Message struct {
	FromPhone   string `json:"fromPhone"`
	ToPhone     string `json:"toPhone"`
	MessageText string `json:"messageText"`
}

//////////////////
//    routes    //
//////////////////

// send whatsapp message
func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// extract details
	var message Message
	json.NewDecoder(r.Body).Decode(&message)

	// send whatsapp message
	err := services.SendWhatsappMessage(message.FromPhone, message.ToPhone, message.MessageText)

	// response - error
	if err != nil {
		panic(err)
	}

	// response - success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"status": true,
			"info":   "Message sent!",
		},
	)
}
