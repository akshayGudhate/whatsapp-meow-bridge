package main

import (
	log "log"
	http "net/http"
	os "os"

	models "akshayGudhate/whatsapp-bridge/src/models"
	routing "akshayGudhate/whatsapp-bridge/src/routes"
	services "akshayGudhate/whatsapp-bridge/src/services"
	handlers "github.com/gorilla/handlers"
)

/////////////////////
//   main method   //
/////////////////////

func main() {
	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	go func() {
		models.StartSyncingToAllExistingDevices()
	}()

	// cors handler
	handlerWithCors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
	)(routing.Routes())

	// het port from env file
	portNumber := services.GetEnvironmentVariables("PORT")

	// server
	srv := &http.Server{
		Addr:     portNumber,
		ErrorLog: errorLog,
		Handler:  handlerWithCors,
	}

	// initialize
	infoLog.Printf("Server is listening on Port URL: http://localhost%s", portNumber)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
