package bridge

//////////////////
//    import    //
//////////////////

import (
	// internal packages
	context "context"
	// external packages
	whatsmeow "go.mau.fi/whatsmeow"
	store "go.mau.fi/whatsmeow/store"
	events "go.mau.fi/whatsmeow/types/events"
	// local packages
	env "akshayGudhate/whatsapp-bridge/src/environment"
)

///////////////////
//   variables   //
///////////////////

var (
	db  Database
	err error
)

//////////////////
//    client    //
//////////////////

// method to establish client connection
func whatsappClientConnection(client *whatsmeow.Client) {
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	if client.Store.ID != nil {
		env.InfoLogger.Println("Connection Established:", client.Store.ID)
	}
}

///////////////////
//  all devices  //
///////////////////

func StartSyncingToAllExistingDevices() {
	// connect to database
	if db.Container == nil {
		db.connectToDatabase()
	}
	// all connected devices updated
	db.getAllConnectedDevices()

	// run goroutines to sync devices
	go func() {
		// connect to all devices
		for _, device := range db.DeviceStore {
			// get client and connect one by one
			meowClient := whatsmeow.NewClient(device, nil)
			// add receive handler
			meowClient.AddEventHandler(eventHandler)
			// connect to client
			whatsappClientConnection(meowClient)
		}
	}()
}

/////////////////////
//  single device  //
/////////////////////

func SyncWithGivenDevice(phone string) string {
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
		if device.ID.User == phone {
			userDevice = device
			break
		}
	}

	// if not add new device
	if userDevice == nil {
		userDevice = db.Container.NewDevice()
		userDevice.Save()
	}

	// create client
	meowClient := whatsmeow.NewClient(userDevice, nil)
	// add receive handler
	meowClient.AddEventHandler(eventHandler)

	//
	// reconnection
	//
	if meowClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := meowClient.GetQRChannel(context.Background())
		// connect to client
		whatsappClientConnection(meowClient)

		for evt := range qrChan {
			if evt.Event == "code" {
				// return the qr code
				return evt.Code
			}
		}

	} else {
		// connect to client
		whatsappClientConnection(meowClient)
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
		receiveMessageEventHandler(v)

	// connection
	case *events.PairSuccess:
		env.InfoLogger.Println("Connection Established:", &v.ID)

	// group info
	case *events.GroupInfo:
		if v.Join != nil {
			env.InfoLogger.Println("Group Details: Group Joined:", v.Join)
		}
		if v.Leave != nil {
			env.InfoLogger.Println("Group Details: Group Leaved:", v.Leave)
		}

		// // default
		// default:
		// 	env.InfoLogger.Println(".........................................")
		// 	env.InfoLogger.Printf("Unknown Event: %+v \n", v)
		// 	env.InfoLogger.Println(".........................................")
	}
}
