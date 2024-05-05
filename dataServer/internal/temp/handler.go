package temp

import (
	"github.com/google/uuid"
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

func Handler(w http.ResponseWriter, r *http.Request) {
	// Handle the request
	log.Println(r.Method, r.URL)
	if r.Method == http.MethodPost {
		createTempObject(w, r)
	} else if r.Method == http.MethodPatch {
		uploadTempObject(w, r)
	} else if r.Method == http.MethodDelete {
		abortTempObject(w, r)
	} else if r.Method == http.MethodPut {
		commitTempObject(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func abortTempObject(w http.ResponseWriter, r *http.Request) {
	uid := strings.Split(r.URL.EscapedPath(), "/")[2]
	i, err := readFromFile(uid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFilePath := filepath.Join(global.Cfg.BaseDir, "/temp/", i.Uuid)
	objectFilePath := infoFilePath + ".dat"
	err = os.Remove(infoFilePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = os.Remove(objectFilePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func commitTempObject(w http.ResponseWriter, r *http.Request) {
	uid := strings.Split(r.URL.EscapedPath(), "/")[2]
	i, err := readFromFile(uid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFilePath := filepath.Join(global.Cfg.BaseDir, "/temp/", i.Uuid)
	objectFilePath := infoFilePath + ".dat"
	err = os.Rename(objectFilePath, filepath.Join(global.Cfg.BaseDir, "/objects/", i.Hash))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	err = os.Remove(infoFilePath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	locate.Add(i.Hash)
	w.WriteHeader(http.StatusOK)
}

func uploadTempObject(w http.ResponseWriter, r *http.Request) {
	uid := strings.Split(r.URL.EscapedPath(), "/")[2]
	i, err := readFromFile(uid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	f, err := os.Create(filepath.Join(global.Cfg.BaseDir, "/temp/", i.Uuid+".dat"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	stat, err := f.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if stat.Size() != i.Size {
		log.Printf("size mismatch: actual %d, exceeds %d\n", stat.Size(), i.Size)
		w.WriteHeader(http.StatusInternalServerError)
		os.Remove(f.Name())
		os.Remove(filepath.Join(global.Cfg.BaseDir, "/temp/", i.Uuid))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createTempObject(w http.ResponseWriter, r *http.Request) {
	size := utils.GetSizeFromHeader(r)
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	uid := uuid.New().String()
	i := info{name, size, uid}
	err := i.writeToFile()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte(uid))
}
