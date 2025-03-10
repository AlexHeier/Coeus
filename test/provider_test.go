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

	resp, err := coeus.Send(coeus.RequestStruct{
		History:      nil,
		Systemprompt: "test",
		Userprompt:   "Hello Ollama"})
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if resp.Response == "" {
		t.Errorf("Expected response, got %v", resp)
	}
	fmt.Printf("Ollama test answer: %v\n", resp.Response)
}
