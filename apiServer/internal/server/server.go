package server

import (
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/locate"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/objects"
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
	go serviceDiscovery.Register(APIServerPrefix)
	go serviceDiscovery.Discovery(APIServerPrefix)
	go serviceDiscovery.Discovery(DataServerPrefix)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatalf(http.ListenAndServe(fmt.Sprintf(":%d", global.Cfg.Port), nil).Error())
}
