package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
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

	hash := strings.Split(r.URL.EscapedPath(), "/")[2]

	f, err := os.Open(filepath.Join(global.Cfg.Server.BaseDir, "objects", hash))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	d := utils.Checksum(f)
	if d != hash {
		log.Println("data corruption detected")
		os.Remove(f.Name())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}
