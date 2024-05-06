package heartbeat

import (
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	"github.com/qiaofufu/tinyoss_kernal/third_party/discovery"
	"math/rand/v2"
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

func ChooseRandomDataServer(n int, exclude map[string]struct{}) []string {
	servers := ServiceDiscovery.GetServers(DataServerPrefix)
	ds := make([]string, 0)
	candidates := make([]string, 0)
	for _, server := range servers {
		if _, ok := exclude[server.(string)]; ok {
			continue
		}
		candidates = append(candidates, server.(string))
	}
	if len(candidates) < n {
		return nil
	}
	p := rand.Perm(len(candidates))
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return ds
}
