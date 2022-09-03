package environment

import (
	// internal packages
	log "log"
	os "os"
)

/////////////////////
//    variables    //
/////////////////////

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

/////////////////////
//   set loggers   //
/////////////////////

func CreateLoggerInstances() {
	// info logs
	InfoLogger = log.New(os.Stdout, "INFO : ", log.Ldate|log.Ltime)
	// error logs
	ErrorLogger = log.New(os.Stderr, "ERROR : ", log.Ldate|log.Ltime|log.Lshortfile)
}
