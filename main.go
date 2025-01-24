package main

import (
	coeus "Coeus/Coeus"

	"log"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	var Coeus coeus.Coeus

	// Setup router with handler functions
	var ServerConfig coeus.SrvConfig
	ServerConfig.Router = mux.NewRouter()
	ServerConfig.HttpProtocol = "http"
	ServerConfig.IPaddress = "localhost"
	ServerConfig.Port = "9050"

	ServerConfig.Router.HandleFunc("/api/v1/status", coeus.StatusHandler)

	LLMConfig := Coeus.NewOllamaConfig("http", "localhost", "11434", "llama3.2", "YEr fucked mate", false)

	Coeus.Init(ServerConfig, LLMConfig)
	Coeus.Start()
}
