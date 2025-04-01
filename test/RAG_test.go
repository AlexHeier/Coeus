package coeus_test

import (
	"os"
	"testing"
	"time"

	coeus "github.com/AlexHeier/Coeus"
)

func TestCoeusRAG(t *testing.T) {
	err := coeus.EnableRAG(os.Getenv("DATABASE_HOST"), "w2v", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), 5432)
	if err != nil {
		t.Fatalf("Failed to enable RAG: %v", err)
	}

	err = coeus.RAGConfig(2, 250, 0.25, 2, "./RAG")
	if err != nil {
		t.Fatalf("Failed to configure RAG: %v", err)
	}

	time.Sleep(5 * time.Minute)

	coeus.GetRAG("What is the deepest oceanic trench, and how deep is it?")
	coeus.GetRAG("Why are bananas considered berries while strawberries are not?")
	coeus.GetRAG("How long is the Great Wall of China, and what was its purpose?")
	coeus.GetRAG("How many hearts does an octopus have, and what is unique about its blood?")
	coeus.GetRAG("Why is lightning hotter than the surface of the sun?")
	coeus.GetRAG("Why doesnt honey spoil, even after thousands of years?")
	coeus.GetRAG("Which has existed longer, sharks or trees?")
	coeus.GetRAG("How does a day on Venus compare to a year on Venus?")
	coeus.GetRAG("Why does the Eiffel Tower grow taller in the summer?")
	coeus.GetRAG("Are there more stars in the universe or grains of sand on Earth?")
}
