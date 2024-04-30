package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	// Handle the request

	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]

	f, err := os.Open(filepath.Join(global.Cfg.Server.BaseDir, "objects", objectName))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
