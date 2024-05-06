package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/third_party/utils"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	// Handle the request

	temp := strings.Split(r.URL.EscapedPath(), "/")[2]
	hash := strings.Split(temp, ".")[0]
	f, err := os.Open(filepath.Join(global.Cfg.Server.BaseDir, "objects", temp))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	d := utils.Checksum(f)
	if d != hash {
		log.Println("data corruption detected")
		locate.Del(hash)
		os.Remove(f.Name())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}
