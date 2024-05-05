package server

import (
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/objects"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/temp"
	"github.com/qiaofufu/tinyoss_kernal/third_party/discovery"
	"log"
	"net/http"
)

const (
	DataServerPrefix = "data-services"
	APIServerPrefix  = "api-services"
)

var serviceDiscovery *discovery.Discovery

func Start() {
	// Start the server
	serviceDiscovery = discovery.NewRecovery(global.Etcd, global.Cfg.Ip, global.Cfg.Port)
	go serviceDiscovery.Register(DataServerPrefix)
	go serviceDiscovery.Discovery(DataServerPrefix)
	go serviceDiscovery.Discovery(APIServerPrefix)
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatalf(http.ListenAndServe(fmt.Sprintf(":%d", global.Cfg.Port), nil).Error())
}
