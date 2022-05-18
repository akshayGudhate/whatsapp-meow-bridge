package services

import (
	proto "google.golang.org/protobuf/proto"

	models "akshayGudhate/whatsapp-bridge/src/models"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	types "go.mau.fi/whatsmeow/types"
)

//////////////////////////
//   whatsapp message   //
//////////////////////////

func SendWhatsappMessage(fromPhone, toPhone, text string) error {
	// encode the data
	recipient, _ := types.ParseJID(toPhone+"@s.whatsapp.net")
	messageText := &waProto.Message{Conversation: proto.String(text)}

	// het client
	if models.MeowClient == nil {
		models.ConnectToGivenClient(fromPhone)
	}

	// send message
	_, err := models.MeowClient.SendMessage(recipient, "", messageText)
	if err != nil {
		return err
	}
	return nil
}
