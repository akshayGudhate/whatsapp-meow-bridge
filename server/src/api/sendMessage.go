package api

import (
	json "encoding/json"
	http "net/http"

	bridge "akshayGudhate/whatsapp-bridge/src/bridge"
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
	responseString := bridge.SendWhatsappMessage(
		message.FromPhone, message.ToPhone, message.MessageText,
	)
	if responseString != "" {
		// response - success
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"status": false,
				"info":   responseString,
			},
		)
		return
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
