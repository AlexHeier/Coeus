package main

import (
	"Coeus/llm"
	"Coeus/provider"
	"fmt"
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

	err := provider.Ollama(os.Getenv("OLLAMA_IP"), os.Getenv("OLLAMA_PORT"), os.Getenv("OLLAMA_MODEL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	con := llm.BeginConversation()
	con.Prompt("Hi im 10 years old and where does babies come from?")

	fmt.Println(con.LatestResponse)

	//fmt.Println(con.DumpConversation())

	con.Prompt("Yes that does make sence! Thank you. Also what was my previous question?")

	//fmt.Println(con.DumpConversation())

	fmt.Println(con.LatestResponse)

	con.Prompt("Do you have any jokes relevant to the last questions i gave you?")

	fmt.Println(con.LatestResponse)

	con.Prompt("Whats PI's first 150 numbers?")

	fmt.Println(con.LatestResponse)

	con.Prompt("Can you give me a list of the previous questions asked?")

	fmt.Println(con.LatestResponse)

	con.Prompt("What part about your functionality do you think im testing?")

	fmt.Println(con.LatestResponse)

	//fmt.Println(con.DumpConversation())

}
