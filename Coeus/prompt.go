package coeus

import "net/http"

func PromptHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}
