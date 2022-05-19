package models

//////////////////
//    import    //
//////////////////

import (
	context "context"
	fmt "fmt"

	_ "github.com/lib/pq"
	whatsmeow "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	store "go.mau.fi/whatsmeow/store"
	sqlstore "go.mau.fi/whatsmeow/store/sqlstore"
	types "go.mau.fi/whatsmeow/types"
	events "go.mau.fi/whatsmeow/types/events"
	proto "google.golang.org/protobuf/proto"

	handlers "akshayGudhate/whatsapp-bridge/src/handlers"
)

//////////////////
//   variable   //
//////////////////

const (
	databaseDialect string = "postgres"
	databaseURL     string = "postgres://akshayg:peshawa8@127.0.0.1:5432/whatsmeow"
)

//////////////////
//    client    //
//////////////////

var MeowClient *whatsmeow.Client

func WhatsappClientConnection(client *whatsmeow.Client) {
	err := client.Connect()
	if err != nil {
		panic(err)
	}
	if client.Store.ID != nil {
		fmt.Println("Connected to --->", client.Store.ID)
	}
}

///////////////////
//  all devices  //
///////////////////

func StartSyncingToAllExistingDevices() {
	// database
	container, err := sqlstore.New(databaseDialect, databaseURL, nil)
	if err != nil {
		panic(err)
	}
	// All devices
	deviceStore, err := container.GetAllDevices()
	if err != nil {
		panic(err)
	}

	go func() {
		// connect to all devices
		for _, device := range deviceStore {
			// get client and connect one by one
			MeowClient = whatsmeow.NewClient(device, nil)
			// add receive handler
			MeowClient.AddEventHandler(eventHandler)
			// connect to client
			WhatsappClientConnection(MeowClient)
		}
	}()
}

/////////////////////
//  single device  //
/////////////////////

func SyncWithGivenDevice(phone string) string {
	// database
	container, err := sqlstore.New(databaseDialect, databaseURL, nil)
	if err != nil {
		panic(err)
	}
	// All devices
	deviceStore, err := container.GetAllDevices()
	if err != nil {
		panic(err)
	}

	// search device and connect
	var userDevice *store.Device
	// check existing devices
	for _, device := range deviceStore {
		if device.ID.User == phone {
			userDevice = device
			break
		}
	}

	// if not add new device
	if userDevice == nil {
		userDevice = container.NewDevice()
	}

	// create client
	MeowClient = whatsmeow.NewClient(userDevice, nil)
	// add receive handler
	MeowClient.AddEventHandler(eventHandler)

	//
	// reconnection
	//
	if MeowClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := MeowClient.GetQRChannel(context.Background())
		// connect to client
		WhatsappClientConnection(MeowClient)

		for evt := range qrChan {
			if evt.Event == "code" {
				// return the qr code
				return evt.Code
			}
		}

	} else {
		// connect to client
		WhatsappClientConnection(MeowClient)
	}
	return ""
}

//////////////////
//    events    //
//////////////////

// Receive Message
func eventHandler(event interface{}) {
	switch v := event.(type) {
	// messages
	case *events.Message:
		handlers.ReceiveMessage(v, MeowClient)

		/** Uncomment below comments to get below info */

		// // group info
		// case *events.GroupInfo:
		// 	if v.Join != nil {
		// 		fmt.Println("Group Details: Group Joined:", v.Join)
		// 	}
		// 	if v.Leave != nil {
		// 		fmt.Println("Group Details: Group Leaved:", v.Leave)
		// 	}

		// default:
		// 	fmt.Println(".........................................")
		// 	fmt.Printf("Unknown Event: %+v \n", v)
		// 	fmt.Println(".........................................")
	}
}

// Send Message
func SendWhatsappMessage(fromPhone, toPhone, text string) string {
	// database
	container, err := sqlstore.New(databaseDialect, databaseURL, nil)
	if err != nil {
		panic(err)
	}
	// All devices
	deviceStore, err := container.GetAllDevices()
	if err != nil {
		panic(err)
	}

	// search device and connect
	var userDevice *store.Device
	// check existing devices
	for _, device := range deviceStore {
		if device.ID.User == fromPhone {
			userDevice = device
			break
		}
	}
	// if not add new device
	if userDevice == nil {
		return "Invalid from phone number"
	}

	// create client
	MeowClient = whatsmeow.NewClient(userDevice, nil)
	// add receive handler
	MeowClient.AddEventHandler(eventHandler)
	// connect to client
	WhatsappClientConnection(MeowClient)

	// encode the data
	recipient, _ := types.ParseJID(toPhone + "@s.whatsapp.net")
	messageText := &waProto.Message{Conversation: proto.String(text)}

	// send message
	_, err = MeowClient.SendMessage(recipient, "", messageText)
	if err != nil {
		return "Something went wrong! Try Again."
	}
	return ""
}
