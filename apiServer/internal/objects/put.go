package objects

import (
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/heartbeat"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objectstream"
	"github.com/qiaofufu/tinyoss_kernal/third_party/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	hash := utils.GetHashFromHeader(r)
	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := r.ContentLength
	log.Println("put", objectName, hash, size)
	err := storeObject(r.Body, url.PathEscape(hash), size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	meta.Add(objectName, hash, size)
	w.WriteHeader(http.StatusOK)
}

func storeObject(r io.Reader, hash string, size int64) error {
	exist, err := locate.Exist(hash)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	stream := putStream(hash, size)
	if stream == nil {
		return fmt.Errorf("no data server available")
	}
	reader := io.TeeReader(r, stream)
	res := utils.Checksum(reader)
	if res != hash {
		err := stream.Abort()
		if err != nil {
			return err
		}
		return fmt.Errorf("checksum failed, actual: %s, expected: %s", res, hash)
	}
	return stream.Commit()
}

func putStream(objectName string, size int64) *objectstream.TempObjectStream {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		log.Println("No data server available")
		return nil
	}
	return objectstream.NewTempObjectStream(server, objectName, size)
}
