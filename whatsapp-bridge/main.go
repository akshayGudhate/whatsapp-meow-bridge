package main

import (
	// internal packages
	http "net/http"
	sync "sync"
	time "time"
	// local packages
	api "akshayGudhate/whatsapp-bridge/src/api"
	bridge "akshayGudhate/whatsapp-bridge/src/bridge"
	env "akshayGudhate/whatsapp-bridge/src/environment"
)

/////////////////////
//   init method   //
/////////////////////

func init() {
	// wait group
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// goroutines for logger initiation
	go env.CreateLoggerInstances(wg)

	// goroutines for syncing client connections
	go bridge.StartSyncingToAllExistingDevices(wg)

	// wait to complete the goroutines execution
	wg.Wait()
}

/////////////////////
//   main method   //
/////////////////////

func main() {
	// server
	app := &http.Server{
		Addr:     env.PORT,
		ErrorLog: env.ErrorLogger,
		Handler:  api.GetAPIRouter(),
		// enforce timeouts for server requests
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// initialize
	env.InfoLogger.Printf("Server is listening on Port URL: http://localhost%s", env.PORT)
	env.ErrorLogger.Fatal(app.ListenAndServe())
}
