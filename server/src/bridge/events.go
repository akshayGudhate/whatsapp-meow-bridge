package bridge

import (
	fmt "fmt"                                // fmt package
	proto "google.golang.org/protobuf/proto" // proto buffers package
	strconv "strconv"                        // type conversion package
	strings "strings"                        // strings package

	whatsmeow "go.mau.fi/whatsmeow"            // bridge - whatsapp - whatsmeow
	waProto "go.mau.fi/whatsmeow/binary/proto" // bridge - whatsapp - binary package
	events "go.mau.fi/whatsmeow/types/events"  // bridge - whatsapp - events packages

	services "akshayGudhate/whatsapp-bridge/src/services" // services local package
)

///////////////////
//   variables   //
///////////////////

// env variables
var (
	TEST_USER1 = services.TEST_USER1
	TEST_USER2 = services.TEST_USER2
)

//////////////////
//    routes    //
//////////////////

// send whatsapp message
func ReceiveMessage(m *events.Message, c *whatsmeow.Client) {
	// only for personal use so remove this condition if you want this for all
	if !(m.Info.Sender.User == TEST_USER1 || m.Info.Sender.User == TEST_USER2) {
		return
	}

	// set message
	var messageText *waProto.Message
	// case
	if *m.Message.Conversation != "" {
		if m.Info.IsGroup {
			// group message
			fmt.Println(
				"Group Details: ", *m.Message.Conversation,
				"---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
				"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
				"---> in group", m.Info.Chat,
			)
		} else {
			// personal message
			fmt.Println(
				"Personal Message: ", *m.Message.Conversation,
				"---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
				"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
			)

			// products details
			products := [...]string{"1: Aditya Kela", "2: Ashwin Achari", "3: Akshay Gudhate", "4: Other"}

			switch strings.ToLower(*m.Message.Conversation) {
			case "hi":
				msg := fmt.Sprintf("Hey, welcome!\nSelect a person you want to contact\n\n%s\n\nEnter option number here. Ex. *2*", strings.Join(products[:], "\n"))
				messageText = &waProto.Message{Conversation: proto.String(msg)}
			case "1", "2", "3", "4":
				idx, _ := strconv.Atoi(*m.Message.Conversation)
				msg := fmt.Sprintf("Thank you for selecting! You have selected *%s*.", strings.Split(products[idx-1], ": ")[1])
				messageText = &waProto.Message{Conversation: proto.String(msg)}
			}

			c.SendMessage(m.Info.Sender, "", messageText)
		}
	} else {
		messageText = &waProto.Message{Conversation: proto.String("Please, enter valid text!")}
		c.SendMessage(m.Info.Sender, "", messageText)
	}
}
