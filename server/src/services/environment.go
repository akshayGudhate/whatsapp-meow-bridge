package services

import (
	godotenv "github.com/joho/godotenv"
)

/////////////////////
//   environment   //
/////////////////////

func GetEnvironmentVariables(variableName string) string {
	myEnvironmentVariable, _ := godotenv.Read()
	return myEnvironmentVariable[variableName]
}
