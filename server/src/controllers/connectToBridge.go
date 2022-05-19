package controllers

import (
	json "encoding/json"
	png "image/png"
	http "net/http"

	barcode "github.com/boombuler/barcode"
	qr "github.com/boombuler/barcode/qr"

	models "akshayGudhate/whatsapp-bridge/src/models"
)

//////////////////
//    routes    //
//////////////////

// send whatsapp message
func GetConnectionQRCode(w http.ResponseWriter, r *http.Request) {
	// parse params
	fromPhone := r.URL.Query().Get("fromPhone")
	if fromPhone != "" {
		// connect to client
		newQRCodeBufferString := models.SyncWithGivenDevice(fromPhone)

		// if not connected then connect
		if newQRCodeBufferString != "" {
			// generate qr code
			qrCode, _ := qr.Encode(newQRCodeBufferString, qr.L, qr.Auto)
			qrCode, _ = barcode.Scale(qrCode, 256, 256)

			// response - qr code
			png.Encode(w, qrCode)

		} else {
			// response - already scanned
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(
				map[string]interface{}{
					"status": false,
					"info":   "This device is already connected!",
				},
			)
		}

	} else {
		// response - already scanned
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			map[string]interface{}{
				"status": false,
				"info":   "Invalid phone number!",
			},
		)
	}
}
