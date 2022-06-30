package main

import (
	log "log"
	http "net/http"
	os "os"
	time "time"

	api "akshayGudhate/whatsapp-bridge/src/api"
	bridge "akshayGudhate/whatsapp-bridge/src/bridge"
	services "akshayGudhate/whatsapp-bridge/src/services"
)

///////////////////
//   variables   //
///////////////////

var PORT = services.PORT

/////////////////////
//   main method   //
/////////////////////

func main() {
	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// goroutines for syncing client connections
	go func() { bridge.StartSyncingToAllExistingDevices() }()

	// server
	srv := &http.Server{
		Addr:     PORT,
		ErrorLog: errorLog,
		Handler:  api.GetHandlerWithRoutes(),
		// enforce timeouts for servers you create
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// initialize
	infoLog.Printf("Server is listening on Port URL: http://localhost%s", PORT)
	errorLog.Fatal(srv.ListenAndServe())
}
