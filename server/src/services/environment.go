package services

import godotenv "github.com/joho/godotenv"

/////////////////////
//   environment   //
/////////////////////

func GetEnvironmentVariables(variableName string) string {
	myEnvironmentVariable, _ := godotenv.Read()
	return myEnvironmentVariable[variableName]
}

/////////////////////
//    variables    //
/////////////////////

var (
	PORT             = GetEnvironmentVariables("PORT")
	DATABASE_DIALECT = GetEnvironmentVariables("DATABASE_DIALECT")
	DATABASE_URL     = GetEnvironmentVariables("DATABASE_URL")
)
