package bridge

import (
	// internal packages
	fmt "fmt"
	strconv "strconv"
	strings "strings"
	// external packages
	whatsmeow "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	events "go.mau.fi/whatsmeow/types/events"
	proto "google.golang.org/protobuf/proto"
	// local packages
	env "akshayGudhate/whatsapp-bridge/src/environment"
)

//////////////////
//    routes    //
//////////////////

// send whatsapp message
func receiveMessage(m *events.Message, c *whatsmeow.Client) {
	// only for personal use so remove this condition if you want this for all
	if !(m.Info.Sender.User == env.TEST_USER1 || m.Info.Sender.User == env.TEST_USER2) {
		return
	}

	// set message
	var messageText *waProto.Message
	// case
	if *m.Message.Conversation != "" {
		if m.Info.IsGroup {
			// group message
			env.InfoLogger.Println(
				"Group Details: ", *m.Message.Conversation,
				"---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
				"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
				"---> in group", m.Info.Chat,
			)
		} else {
			// personal message
			env.InfoLogger.Println(
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
