package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objectstream"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionID := r.URL.Query().Get("version")
	var (
		version int
		err     error
	)
	if len(versionID) > 0 {
		version, err = strconv.Atoi(versionID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return

		}
	}
	objectMeta, err := meta.GetMetadata(objectName, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if objectMeta.Hash == "" {
		log.Println("Object already deleted")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream, err := getStream(url.PathEscape(objectMeta.Hash))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = io.Copy(w, stream)
}

func getStream(objectName string) (io.Reader, error) {
	server, err := locate.LocateFromAllServer(objectName)
	if err != nil {
		return nil, err
	}
	return objectstream.NewGetStream(server, objectName)
}
