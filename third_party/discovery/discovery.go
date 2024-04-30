package discovery

import (
	"context"
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

type Discovery struct {
	ip          string
	port        int32
	servers     map[string]map[string]any
	serverMutex sync.Mutex
	etcd        *etcdClientv3.Client
}

func NewRecovery(etcd *etcdClientv3.Client, ip string, port int32) *Discovery {
	return &Discovery{
		servers: make(map[string]map[string]any),
		etcd:    etcd,
		ip:      ip,
		port:    port,
	}
}

func (r *Discovery) Discovery(prefix string) {
	serverChan := r.etcd.Watch(context.Background(), fmt.Sprintf("/%s/", prefix), etcdClientv3.WithPrefix())
	for resp := range serverChan {
		for _, ev := range resp.Events {
			switch ev.Type {
			case etcdClientv3.EventTypePut:
				r.addServer(prefix, string(ev.Kv.Key), string(ev.Kv.Value))
			case etcdClientv3.EventTypeDelete:
				r.removeServer(prefix, string(ev.Kv.Key))
			}
		}
	}
}

func (r *Discovery) Register(prefix string) {
	// Grant lease
	lease, err := r.etcd.Grant(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	}
	// KeepAlive lease
	keepaliveCh, err := r.etcd.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		log.Fatal(err)
	}
	// Register service
	key := fmt.Sprintf("/%s/%s:%d", prefix, r.ip, r.port)
	value := fmt.Sprintf("%s:%d", r.ip, r.port)
	put, err := global.Etcd.Put(context.Background(), key, value, etcdClientv3.WithLease(lease.ID))
	if err != nil {
		log.Fatal(put)
		return
	}

	for {
		select {
		case _, ok := <-keepaliveCh:
			if !ok {
				log.Fatal("keepalive channel closed")
			}
		}
	}
}

func (r *Discovery) addServer(prefix, key, value string) {
	r.serverMutex.Lock()
	defer r.serverMutex.Unlock()
	if _, ok := r.servers[key]; !ok {
		r.servers[key] = make(map[string]any)
	}
	r.servers[prefix][key] = value
}

func (r *Discovery) removeServer(prefix, key string) {
	r.serverMutex.Lock()
	defer r.serverMutex.Unlock()
	if _, ok := r.servers[prefix]; ok {
		delete(r.servers[prefix], key)
		if len(r.servers[prefix]) == 0 {
			delete(r.servers, prefix)
		}
	}
}
