package main

import (
	log "log"
	http "net/http"
	os "os"

	api "akshayGudhate/whatsapp-bridge/src/api"
	bridge "akshayGudhate/whatsapp-bridge/src/bridge"
)

/////////////////////
//   main method   //
/////////////////////

func main() {
	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	go func() {
		// goroutines for syncing client connections
		bridge.StartSyncingToAllExistingDevices()
	}()

	// het port from env file
	portNumber := bridge.GetEnvironmentVariables("PORT")

	// server
	srv := &http.Server{
		Addr:     portNumber,
		ErrorLog: errorLog,
		Handler:  api.GetHandlerWithRoutes(),
	}

	// initialize
	infoLog.Printf("Server is listening on Port URL: http://localhost%s", portNumber)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
