package models

//////////////////
//    import    //
//////////////////

import (
	context "context"
	fmt "fmt"

	_ "github.com/lib/pq"
	whatsmeow "go.mau.fi/whatsmeow"
	sqlstore "go.mau.fi/whatsmeow/store/sqlstore"
	types "go.mau.fi/whatsmeow/types"
	events "go.mau.fi/whatsmeow/types/events"

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

///////////////////
//  all devices  //
///////////////////

func ConnectToAllClients() {
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

	//
	// connect to all devices
	//
	fmt.Println("------------------------ Connecting to existing devices ------------------------")
	for i := 0; i < len(deviceStore); i++ {
		// get client and connect one by one
		MeowClient = whatsmeow.NewClient(deviceStore[i], nil)
		// add receive handler
		MeowClient.AddEventHandler(eventHandler)
		err = MeowClient.Connect()
		fmt.Println("Connected to client --->", MeowClient.Store.ID)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("------------------------ Connected to all the devices ------------------------")
}

/////////////////////
//  single device  //
/////////////////////

func ConnectToGivenClient(phone string) string {
	// database
	container, err := sqlstore.New(databaseDialect, databaseURL, nil)
	if err != nil {
		panic(err)
	}

	// TODO:
	// 1. get phone number from the DB
	// 2. get JID and parse of user for connection
	// 2. check weather device already present for JID
	// 3. if device available then connect
	// 4. if device not available then create device & then connect

	// search device and connect
	singleDevice := container.NewDevice()

	// create client
	MeowClient = whatsmeow.NewClient(singleDevice, nil)
	// add receive handler
	MeowClient.AddEventHandler(eventHandler)
	newClient, _ := types.ParseJID("919561214185@s.whatsapp.net")

	//
	// reconnection
	//
	if MeowClient.Store.ID == nil || *MeowClient.Store.ID != newClient {
		// No ID stored, new login
		qrChan, _ := MeowClient.GetQRChannel(context.Background())
		err = MeowClient.Connect()
		fmt.Println("Connected to client --->", MeowClient.Store.ID)
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				// return the qr code
				return evt.Code
			}
		}

	} else {
		// Already logged in, just connect
		err = MeowClient.Connect()
		fmt.Println("Connected to client --->", MeowClient.Store.ID)
		if err != nil {
			panic(err)
		}
	}
	return ""
}

//////////////////
//    events    //
//////////////////

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
