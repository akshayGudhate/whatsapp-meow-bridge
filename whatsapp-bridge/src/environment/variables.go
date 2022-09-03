package environment

// external packages
import godotenv "github.com/joho/godotenv"

/////////////////////
//  get variables  //
/////////////////////

func getEnvironmentVariables(variableName string) string {
	myEnvironmentVariable, _ := godotenv.Read()
	return myEnvironmentVariable[variableName]
}

/////////////////////
//  set variables  //
/////////////////////

var (
	PORT             = getEnvironmentVariables("PORT")
	DATABASE_URL     = getEnvironmentVariables("DATABASE_URL")
	DATABASE_DIALECT = getEnvironmentVariables("DATABASE_DIALECT")
	TEST_USER1       = getEnvironmentVariables("TEST_USER1")
	TEST_USER2       = getEnvironmentVariables("TEST_USER2")
	TEST_USER3       = getEnvironmentVariables("TEST_USER3")
	TEST_USER4       = getEnvironmentVariables("TEST_USER4")
)
