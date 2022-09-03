package bridge

//////////////////
//    import    //
//////////////////

import (
	// internal packages
	context "context"
	// external packages
	_ "github.com/lib/pq"
	whatsmeow "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	store "go.mau.fi/whatsmeow/store"
	sqlstore "go.mau.fi/whatsmeow/store/sqlstore"
	types "go.mau.fi/whatsmeow/types"
	events "go.mau.fi/whatsmeow/types/events"
	proto "google.golang.org/protobuf/proto"
	// local packages
	env "akshayGudhate/whatsapp-bridge/src/environment"
)

///////////////////
//   variables   //
///////////////////

// database struct
type Database struct {
	Container   *sqlstore.Container
	DeviceStore []*store.Device
}

// placeholders
var (
	err        error
	db         Database
	MeowClient *whatsmeow.Client
)

//////////////////
//   database   //
//////////////////

// method to connect database
func (db *Database) connectToDatabase() {
	// connection
	db.Container, err = sqlstore.New(env.DATABASE_DIALECT, env.DATABASE_URL, nil)
	if err != nil {
		panic(err)
	}
}

// method to get all connect devices
func (db *Database) getAllConnectedDevices() {
	db.DeviceStore, err = db.Container.GetAllDevices()
	if err != nil {
		panic(err)
	}
}

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
		env.InfoLogger.Println("Connected to --->", client.Store.ID)
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
			MeowClient = whatsmeow.NewClient(device, nil)
			// add receive handler
			MeowClient.AddEventHandler(eventHandler)
			// connect to client
			whatsappClientConnection(MeowClient)
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
		whatsappClientConnection(MeowClient)

		for evt := range qrChan {
			if evt.Event == "code" {
				// return the qr code
				return evt.Code
			}
		}

	} else {
		// connect to client
		whatsappClientConnection(MeowClient)
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
		receiveMessage(v, MeowClient)
	case *events.PairSuccess:
		env.InfoLogger.Println("Connected to --->", &v.ID)

		/** Uncomment below comments to get below info */

		// // group info
		// case *events.GroupInfo:
		// 	if v.Join != nil {
		// 		env.InfoLogger.Println("Group Details: Group Joined:", v.Join)
		// 	}
		// 	if v.Leave != nil {
		// 		env.InfoLogger.Println("Group Details: Group Leaved:", v.Leave)
		// 	}

		// default:
		// 	env.InfoLogger.Println(".........................................")
		// 	fmt.Printf("Unknown Event: %+v \n", v)
		// 	env.InfoLogger.Println(".........................................")
	}
}

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
	MeowClient = whatsmeow.NewClient(userDevice, nil)
	// add receive handler
	MeowClient.AddEventHandler(eventHandler)
	// connect to client
	whatsappClientConnection(MeowClient)

	// encode the data
	recipient, _ := types.ParseJID(toPhone + "@s.whatsapp.net")
	messageText := &waProto.Message{Conversation: proto.String(text)}

	//
	// send message
	//
	go func() { _, err = MeowClient.SendMessage(recipient, "", messageText) }()
	if err != nil {
		return "Something went wrong! Try Again."
	}
	return ""
}
