package objects

import (
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	log.Println(r.Method, r.URL)
	if r.Method == http.MethodPut {
		put(w, r)
	} else if r.Method == http.MethodGet {
		get(w, r)
	} else if r.Method == http.MethodDelete {
		del(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
