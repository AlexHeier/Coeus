package main

import (
	"Coeus/conversation"
	"Coeus/conversation/memory"
	"Coeus/provider"
	"fmt"
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

	llm, err := provider.Ollama(ip, port, model)
	if err != nil {
		log.Fatal(err.Error())
	}

	var con conversation.Struct
	con.Setup(llm, memory.Summary)

	fmt.Print(con)
}
