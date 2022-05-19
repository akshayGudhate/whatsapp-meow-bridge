package controllers

import (
	json "encoding/json"
	http "net/http"

	models "akshayGudhate/whatsapp-bridge/src/models"
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
	responseString := models.SendWhatsappMessage(message.FromPhone, message.ToPhone, message.MessageText)
	if responseString != "" {
		// response - success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"status": false,
				"info":   responseString,
			},
		)
	} else {
		// response - success
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"status": true,
				"info":   "Message sent!",
			},
		)
	}
}
