package versions

import (
	"encoding/json"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]

	metas, err := meta.SearchAllVersions(objectName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(metas)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
