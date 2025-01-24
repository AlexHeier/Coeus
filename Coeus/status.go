package coeus

import "net/http"

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		statusGetHandler(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

func statusGetHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "This shows the endpoint works :) "+r.Method, http.StatusOK)
}
