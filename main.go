package main

import (
	"Coeus/dashboard"
	"Coeus/llm"
	"Coeus/llm/tool"
	"Coeus/provider"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	err := provider.Azure(os.Getenv("AZURE_ENDPOINT"), os.Getenv("AZURE_API_KEY"), 1.0, 120)
	//err := provider.OpenAI("gpt-4", os.Getenv("OPENAI_API_KEY"))
	//err := provider.Ollama(os.Getenv("OLLAMA_IP"), os.Getenv("OLLAMA_PORT"), os.Getenv("OLLAMA_MODEL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	llm.SetPersona("You are a chatbot with some tools available. Use these when and only when necessary. Be a good bot :)")

	llm.MemoryVersion(llm.MemoryAllMessage)

	//tool.New("Multiply", "Takes two ints and returns the multiplied result.", Multiply)
	//tool.New("GetCurrentTime", "Gets the current time.", GetCurrentTime)
	tool.New("GetMagicData", "Retreives the magic data.", GetMagicData)

	go TimeOutConversations()

	dashboard.Start("9002")
}

func GetCurrentTime() (string, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}

func Multiply(a, b int) int {

	fmt.Printf("Issued Command: Multiply %v by %v\n", a, b)
	return a * b
}

func GetMagicData() (string, error) {
	return "69420", nil
}

func TimeOutConversations() {

	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `INSERT INTO conversations (history) VALUES ($1)`

	for {
		time.Sleep(1 * time.Minute)
		var temp []*llm.Conversation
		for _, c := range llm.ConvAll.Conversations {
			c.Mutex.Lock()
			if time.Since(c.LastActive) > 10*time.Minute {
				_, err := db.Exec(query, c.DumpConversation())
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				temp = append(temp, c)
			}
			c.Mutex.Unlock()
		}
		llm.ConvAll.Conversations = temp
	}
}
