package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"log"
	"net/http"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	log.Println(r.Method, r.URL)
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	err := meta.Del(objectName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
