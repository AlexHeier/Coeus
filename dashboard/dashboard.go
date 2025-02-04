package dashboard

import (
	"fmt"
	"net/http"
	"os"
)

// Enables the built-in dashboard for Coeus. Is usually disabled unless this function is called with a port as arg
func Enable(Port string) error {
	http.HandleFunc("/", webHandler)
	return http.ListenAndServe(":9002", nil)
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
	data, err := os.ReadFile("./dashboard/dashboardv2.html")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	w.Write(data)
}
