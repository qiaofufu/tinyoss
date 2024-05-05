package server

import (
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/heartbeat"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objects"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/versions"
	"log"
	"net/http"
)

func Start() {
	// Start the server
	heartbeat.Init()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatalf(http.ListenAndServe(fmt.Sprintf(":%d", global.Cfg.Port), nil).Error())
}
