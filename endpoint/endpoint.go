package endpoint

import (
	"Coeus/llm"
	"fmt"
	"net/http"
)

const API_VERSION = "v1"
const BASE_URL = "/api/" + API_VERSION
const ENDPOINT_CHAT = BASE_URL + "/chat"
const ENDPOINT_STATUS = BASE_URL + "/status"

//const ENDPOINT_CONVERSATION = BASE_URL + "/conversation"

func ChatHandler(w http.ResponseWriter, r *http.Request, con *llm.Conversation) {
	switch r.Method {
	case http.MethodPost:
		chatPostHandler(w, r, con)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func chatPostHandler(w http.ResponseWriter, r *http.Request, con *llm.Conversation) {

	if con == nil {
		fmt.Println("Conversation cannot be empty")
		http.Error(w, "Missing conversation struct", http.StatusInternalServerError)
		return
	}

}
