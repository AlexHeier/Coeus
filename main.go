package main

import (
	"Coeus/llm"
	"Coeus/llm/memory"
	"Coeus/provider"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	prov, err := provider.Ollama(os.Getenv("OLLAMA_IP"), os.Getenv("OLLAMA_PORT"), os.Getenv("OLLAMA_MODEL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	llm.Setup(prov, memory.Summary)

}
