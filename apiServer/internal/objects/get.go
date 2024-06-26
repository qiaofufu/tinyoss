package objects

import (
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/heartbeat"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/meta"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/rs"
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
	stream, err := getStream(url.PathEscape(objectMeta.Hash), objectMeta.Size)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func getStream(hash string, size int64) (io.Reader, error) {
	servers, err := locate.LocateFromAllServer(hash)
	if err != nil {
		return nil, err
	}
	if len(servers) < global.Cfg.RS.DataShard {
		return nil, fmt.Errorf("can't find enough data server")
	}
	excludeServers := make(map[string]struct{})
	for i := range servers {
		excludeServers[servers[i]] = struct{}{}
	}
	ds := make([]string, 0)
	if len(servers) < global.Cfg.RS.ShardAllNum {
		ds = heartbeat.ChooseRandomDataServer(global.Cfg.RS.ShardAllNum-len(servers), excludeServers)
	}

	return rs.NewGetStream(servers, ds, hash, size)
}
