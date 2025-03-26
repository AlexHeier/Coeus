package coeus_test

import (
	"fmt"
	"os"
	"testing"

	coeus "github.com/AlexHeier/Coeus"
)

func TestOllama(t *testing.T) {

	result := coeus.Ollama("a", "8080", "test")
	if result == nil {
		t.Errorf("Expected err, got %v", result)
	}

	result = coeus.Ollama("10.0.0.1", "s", "test")
	if result == nil {
		t.Errorf("Expected err, got %v", result)
	}

	result = coeus.Ollama("10.0.0.1", "8080", "")
	if result == nil {
		t.Errorf("Expected err, got %v", result)
	}

	result = coeus.Ollama(os.Getenv("OLLAMA_IP"), os.Getenv("OLLAMA_PORT"), os.Getenv("OLLAMA_MODEL"))
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}

	con := coeus.BeginConversation()
	if con == nil {
		t.Errorf("Expected new conversation pointer, got %v", con)
	}

	resp, err := con.Prompt("Hello Ollama")
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}

	if resp.Response == "" {
		t.Errorf("Expected response, got %v", resp)
	}
	fmt.Printf("Ollama test answer: %v\n", resp.Response)
}
