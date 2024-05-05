package heartbeat

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/third_party/discovery"
)

const (
	DataServerPrefix = "data-services"
	APIServerPrefix  = "api-services"
)

var ServiceDiscovery *discovery.Discovery

func Init() {
	ServiceDiscovery = discovery.NewRecovery(global.Etcd, global.Cfg.Ip, global.Cfg.Port)
	go ServiceDiscovery.Register(APIServerPrefix)
	go ServiceDiscovery.Discovery(APIServerPrefix)
	go ServiceDiscovery.Discovery(DataServerPrefix)
}

func ChooseRandomDataServer() string {
	servers := ServiceDiscovery.GetServers(DataServerPrefix)
	for _, server := range servers {
		return server.(string)
	}
	return ""
}
