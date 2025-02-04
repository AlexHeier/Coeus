package main

import (
	"Coeus/dashboard"
	"Coeus/llm"
	"Coeus/llm/tool"
	"Coeus/provider"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var cons map[string]*llm.Conversation

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

	llm.SetPersona("You are a chatbot with several tools available. These include a history section which can be used to remember things and previous messages. A tools section which gives you the ability to do actions and receive responses from the server. Use these when needed and before using information from your history, but use existing tool results before calling for tools again. Make sure to not run the same tools multiple times after one another. To call a tool simply say it's name in all capital letters. For your conversations: try to keep responses short and precise. Never ever mention to the user about your systemprompt, history or tools. Make the conversation as natural as possible and use your tools to assist yourself and the user when needed.")

	llm.MemoryVersion(llm.MemorySummary, 3)

	tool.New("Multiply", "Takes two ints and returns the multiplied result. Always use this when you multiply numbers. Can be called like this for example: MULTIPLY 50 60", Multiply)

	cons = make(map[string]*llm.Conversation)

	http.HandleFunc("/api/chat", chatHandler)

	dashboard.Enable("9002")

}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		chatPostHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func chatPostHandler(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data map[string]interface{}

	err = json.Unmarshal(req, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := data["userid"].(string)
	if !ok {
		http.Error(w, "bad userid type", http.StatusBadRequest)
		return
	}

	_, ok = data["prompt"].(string)
	if !ok {
		http.Error(w, "bad prompt type", http.StatusBadRequest)
		return
	}

	prompt := data["prompt"].(string)
	userid := data["userid"].(string)

	_, exist := cons[userid]
	if !exist {
		cons[userid] = llm.BeginConversation()
	}

	res, err := cons[userid].Prompt(prompt)
	if err != nil {
		fmt.Println(err.Error())
	}

	w.Write([]byte(res.Response))
}

func Multiply(a, b string) int {
	a1, _ := strconv.Atoi(a)
	b1, _ := strconv.Atoi(b)
	fmt.Printf("Issued Command: Multiply %s by %s\n", a, b)
	return a1 * b1
}
