package main

import (
	"Coeus/llm/memory"
	"Coeus/llm/tool"
	"Coeus/provider"
	"fmt"
	"log"
	"os"
)

func main() {

	err := provider.Ollama(os.Getenv("OLLAMA_IP"), os.Getenv("OLLAMA_PORT"), os.Getenv("OLLAMA_MODEL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	memory.Version(memory.Summary)

	tool.New("test", "a test function", test)

	t, err := tool.Find("test")
	if err != nil {
		log.Fatal(err.Error())
	}

	anw, err := t.Run(40, 60)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(anw)

}

func test(a, b int) int {

	return a * b
}
