package bridge

import (
	// external packages
	whatsmeow "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	store "go.mau.fi/whatsmeow/store"
	types "go.mau.fi/whatsmeow/types"
	proto "google.golang.org/protobuf/proto"
)

// Send Message
func SendWhatsappMessage(fromPhone, toPhone, text string) string {
	// connect to database
	if db.Container == nil {
		db.connectToDatabase()
	}
	// all connected devices updated
	db.getAllConnectedDevices()

	// search device and connect
	var userDevice *store.Device
	// check existing devices
	for _, device := range db.DeviceStore {
		if device.ID.User == fromPhone {
			userDevice = device
			break
		}
	}
	// if not add new device
	if userDevice == nil {
		return "Invalid from phone number."
	}

	// create client
	meowClient := whatsmeow.NewClient(userDevice, nil)
	// add receive handler
	meowClient.AddEventHandler(eventHandler)
	// connect to client
	whatsappClientConnection(meowClient)

	// encode the data
	recipient, _ := types.ParseJID(toPhone + "@s.whatsapp.net")
	messageText := &waProto.Message{Conversation: proto.String(text)}

	//
	// send message
	//
	go func() { _, err = meowClient.SendMessage(recipient, "", messageText) }()
	if err != nil {
		return "Something went wrong! Try Again."
	}
	return ""
}
