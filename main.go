package main

import (
	"Coeus/conversation"
	_ "Coeus/conversation/memory"
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

	llm, err := provider.NewOllama(ip, port, model)
	if err != nil {
		log.Fatal(err.Error())
	}

	conversation := conversation.ConversationSetup{
		llm:    llm,
		Memory: conversation.Memory.Summery(),
	}

	fmt.Print(conversation)
}
