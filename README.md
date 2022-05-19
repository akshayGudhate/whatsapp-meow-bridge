# whatsmeow-whatsapp-bridge
- By Akshay Gudhate

whatsmeow is a Go library for the WhatsApp web multidevice API created by Tulir Asokan.
[![Go Reference](https://pkg.go.dev/badge/go.mau.fi/whatsmeow.svg)](https://pkg.go.dev/go.mau.fi/whatsmeow)


## Project Features
Using the above library we I have created a HTTP server, which can handle multiple `WhatsApp*` devices simultaneously.

* Receiving all messages and save them for future use.
* Receiving Group messages and join leave notification.
* Automated replies + tried a bot.
* Message sending over HTTP call.
* New device connection over HTTP call.
* Whatsapp contact list.

### How to use

- Clone the project using `https://github.com/akshayGudhate/whatsapp-meow-bridge.git`.
- Open server folder `cd server` where the HTTPS server code has been added.
- Run `go mod tidy` in server folder.
- Run `go build` and then `./whatsapp-bridge` to run the code or simply Run `go run .` to start the HTTP server.

### Supported API's
- Adding `API-Collection.json` file for the reference

* 1. Connect to server: `http://localhost:8080/api/connect/qr?fromPhone=9195xxxxxx85` provide *fromPhone* --> Scan the QR code using whatsapp.
* 2. Send whatsapp message: `http://localhost:8080/api/whatsapp/send` provide *fromPhone*, *toPhone* and *messageText*.
