package endpoint

import (
	"Coeus/conversation"
	"fmt"
	"net/http"
)

const API_VERSION = "v1"
const BASE_URL = "/api/" + API_VERSION
const ENDPOINT_CHAT = BASE_URL + "/chat"
const ENDPOINT_STATUS = BASE_URL + "/status"

//const ENDPOINT_CONVERSATION = BASE_URL + "/conversation"

func chatHandler(w http.ResponseWriter, r *http.Request, con *conversation.Struct) {
	switch r.Method {
	case http.MethodPost:
		chatPostHandler(w, r, con)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func chatPostHandler(w http.ResponseWriter, r *http.Request, con *conversation.Struct) {

	if con == nil {
		fmt.Println("Conversation cannot be empty")
		http.Error(w, "Missing conversation struct", http.StatusInternalServerError)
		return
	}

}
