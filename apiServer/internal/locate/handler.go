package locate

import (
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	log.Println(r.Method, r.URL)
	if r.Method == http.MethodGet {
		get(w, r)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server, err := LocateFromAllServer(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if server == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = w.Write([]byte(server))
	if err != nil {
		log.Println(err)
	}
}
