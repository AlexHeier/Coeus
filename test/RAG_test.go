package coeus_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	coeus "github.com/AlexHeier/Coeus"
)

func TestCoeusRAG(t *testing.T) {

	err := coeus.RAGConfig(2, 1000, 0.25, 2)
	if err != nil {
		t.Fatalf("Failed to configure RAG: %v", err)
	}

	err = coeus.EnableRAG(os.Getenv("DATABASE_HOST"), "w2v", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), 5432, "./RAG")
	if err != nil {
		t.Fatalf("Failed to enable RAG: %v", err)
	}

	time.Sleep(181 * time.Second)

	line := "\n-------------------------------------------------------------------\n"

	print("What is the name of the store\n")
	fmt.Println(coeus.GetRAG("What is the name of the store"))
	print(line)

	print("What payment methods do you accept\n")
	fmt.Println(coeus.GetRAG("What payment methods do you accept"))
	print(line)

	print("When are you open today?\n")
	fmt.Println(coeus.GetRAG("When are you open today?"))
	print(line)

	print("What store locations do you have\n")
	fmt.Println(coeus.GetRAG("What store locations do you have"))
	print(line)
}
