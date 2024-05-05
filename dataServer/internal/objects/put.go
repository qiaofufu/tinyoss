package objects

import (
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// Handle the request

	dir := filepath.Join(global.Cfg.Server.BaseDir, "objects")
	objectName := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, err := os.Create(filepath.Join(dir, objectName))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
}
