package dashboard

import (
	"Coeus/llm"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var cons map[string]*llm.Conversation

// Enables the built-in dashboard for Coeus. Is usually disabled unless this function is called with a port as arg
func Enable(Port string) error {
	cons = make(map[string]*llm.Conversation)
	http.HandleFunc("/api/chat", chatHandler)
	http.HandleFunc("/", webHandler)
	return http.ListenAndServe(":"+Port, nil)
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
	data, err := os.ReadFile("./dashboard/index.html")
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

	w.Write([]byte(res.Response))
}
