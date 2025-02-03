package main

import (
	"Coeus/llm"
	"Coeus/provider"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

	llm.MemoryVersion(llm.MemoryAllMessage)

	cons = make(map[string]*llm.Conversation)

	http.HandleFunc("/", webHandler)
	http.HandleFunc("/api/chat", chatHandler)
	http.ListenAndServe(":9002", nil)

}

func webHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		webGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func webGetHandler(w http.ResponseWriter, r *http.Request) {
	_ = r
	data, err := os.ReadFile("./index.html")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.Write(data)
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

	w.Write([]byte(res["response"].(string)))
}
