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
	"strconv"
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

	err := provider.Azure(os.Getenv("AZURE_ENDPOINT"), os.Getenv("AZURE_API_KEY"), 1.0, 16)
	if err != nil {
		log.Fatal(err.Error())
	}

	llm.SetPersona("Respond in the language of the last user message. You are a chatbot with tools for memory and actions. Use them when needed, and dont call them multiple times. Keep responses short and natural.")

	llm.MemoryVersion(llm.MemoryAllMessage)

	//tool.New("Multiply", "Takes two ints and returns the multiplied result. ALWAYS call this when multiplying two numbers.", Multiply)
	tool.New("GetMagicData", "Retreives the magic data.", GetMagicData)

	tool.New("GetCurrentTime", "Gets the current time.", GetCurrentTime)

	//go TimeOutConversations()

	dashboard.Start("9002")

}

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetMagicData() string {
	return "69420"
}

func Multiply(a, b string) int {

	fmt.Printf("Issued Command: Multiply %s by %s\n", a, b)

	a1, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}

	b1, err := strconv.Atoi(b)
	if err != nil {
		return 0
	}

	return a1 * b1
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
