package bridge

import (
	// internal packages
	fmt "fmt"
	strconv "strconv"
	strings "strings"
	// external packages
	phonenumbers "github.com/nyaruka/phonenumbers"
	events "go.mau.fi/whatsmeow/types/events"
	// local packages
	env "akshayGudhate/whatsapp-bridge/src/environment"
)

/////////////////////
//  handle events  //
/////////////////////

// send whatsapp message
func receiveMessageEventHandler(m *events.Message, eventReceivedPhone string) {
	// parse country code
	userPhone, _ := phonenumbers.Parse("+"+m.Info.Sender.User, "")
	userCountryCode := phonenumbers.GetCountryCodeForRegion(phonenumbers.GetRegionCodeForNumber(userPhone))

	// simple logger
	if m.Info.IsGroup {
		// group message
		env.InfoLogger.Println(
			"Group Details: ", *m.Message.Conversation,
			// "---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
			// "---> to: ", "(", *m.RawMessage.DeviceSentMessage.DestinationJid, ")",
			"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
			"---> in group", m.Info.Chat,
			"---> country code", userCountryCode,
		)

	} else {
		// personal message
		env.InfoLogger.Println(
			"Personal Message: ", *m.Message.Conversation,
			// "---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
			// "---> to: ", "(", *m.RawMessage.DeviceSentMessage.DestinationJid, ")",
			"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
			"---> country code", userCountryCode,
		)
	}

	// only for personal use so remove this condition if you want this for all
	// if !(m.Info.Sender.User == env.TEST_USER1 || m.Info.Sender.User == env.TEST_USER2) {
	if !(m.Info.Sender.User == env.TEST_USER4) && *m.Message.Conversation == "" {
		return
	}

	// case
	if *m.Message.Conversation != "" {
		if m.Info.IsGroup {
			// group message
			env.InfoLogger.Println(
				"Group Details: ", *m.Message.Conversation,
				// "---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
				// "---> to: ", "(", *m.RawMessage.DeviceSentMessage.DestinationJid, ")",
				"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
				"---> in group", m.Info.Chat,
				"---> country code", userCountryCode,
			)

		} else {
			// personal message
			env.InfoLogger.Println(
				"Personal Message: ", *m.Message.Conversation,
				// "---> to: ", "(", m.Info.DeviceSentMeta.DestinationJID, ")",
				// "---> to: ", "(", *m.RawMessage.DeviceSentMessage.DestinationJid, ")",
				"---> from: ", m.Info.PushName, "(", m.Info.Sender, ")",
				"---> country code", userCountryCode,
			)

			// products details
			products := [...]string{"1: Aditya Kela", "2: Ashwin Achari", "3: Akshay Gudhate", "4: Other"}

			switch strings.ToLower(*m.Message.Conversation) {
			case "hi":
				messageText := fmt.Sprintf("Hey, welcome!\nSelect a person you want to contact\n\n%s\n\nEnter option number here. Ex. *2*", strings.Join(products[:], "\n"))
				fromPhone := strings.Split(m.Info.DeviceSentMeta.DestinationJID, "@")[0]
				// send message
				env.InfoLogger.Println(fromPhone, m.Info.Sender.User, messageText)
				SendWhatsappMessage(&fromPhone, &m.Info.Sender.User, &messageText)
			case "1", "2", "3", "4":
				idx, _ := strconv.Atoi(*m.Message.Conversation)
				messageText := fmt.Sprintf("Thank you for selecting! You have selected *%s*.", strings.Split(products[idx-1], ": ")[1])
				fromPhone := strings.Split(m.Info.DeviceSentMeta.DestinationJID, "@")[0]
				// send message
				env.InfoLogger.Println(fromPhone, m.Info.Sender.User, messageText)
				SendWhatsappMessage(&fromPhone, &m.Info.Sender.User, &messageText)
			}
		}

	} else {
		messageText := "Please, enter valid text!"
		fromPhone := strings.Split(m.Info.DeviceSentMeta.DestinationJID, "@")[0]
		// send message
		env.InfoLogger.Println(fromPhone, m.Info.Sender.User, messageText)
		SendWhatsappMessage(&fromPhone, &m.Info.Sender.User, &messageText)
	}
}
