package coeus_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// Calls all test functions in the package
// go test -v ./test/.

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	code := m.Run()
	os.Exit(code)
}
