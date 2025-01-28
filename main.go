package main

import (
	"Coeus/llm"
	"Coeus/llm/memory"
	"Coeus/provider"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

var ip = "10.212.168.203"
var port = "11434"
var model = "llama3.2"

func main() {

	prov, err := provider.Ollama(ip, port, model)
	if err != nil {
		log.Fatal(err.Error())
	}

	llm.Setup(prov, memory.Summary)

}
