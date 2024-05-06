package locate

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	log.Println(r.Method, r.URL)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info, err := LocateFromAllServer(strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(info)
	w.WriteHeader(http.StatusOK)
}
