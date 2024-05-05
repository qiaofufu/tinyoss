package locate

import (
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
	if info == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, _ = w.Write([]byte(info))
	w.WriteHeader(http.StatusOK)
}
