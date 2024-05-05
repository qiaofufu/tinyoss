package discovery

import (
	"context"
	"fmt"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"time"
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
	r.serverMutex.Lock()
	if _, ok := r.servers[prefix]; !ok {
		r.servers[prefix] = make(map[string]any)
	}
	r.serverMutex.Unlock()
	log.Printf("discovery prefix: %s, key: %s", prefix, fmt.Sprintf("/%s/", prefix))
	serverChan := r.etcd.Watch(context.Background(), fmt.Sprintf("/%s/", prefix), etcdClientv3.WithPrefix())
	resp, err := r.etcd.Get(context.Background(), fmt.Sprintf("/%s/", prefix), etcdClientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return

	}
	for _, kv := range resp.Kvs {
		r.addServer(prefix, string(kv.Key), string(kv.Value))
	}
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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Grant lease
	lease, err := r.etcd.Grant(ctx, 5)
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
	log.Printf("register key: %s, value: %s", key, value)
	put, err := r.etcd.Put(context.Background(), key, value, etcdClientv3.WithLease(lease.ID))
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
	log.Printf("add server prefix: %s, key: %s, value: %s", prefix, key, value)
	r.servers[prefix][key] = value

}

func (r *Discovery) removeServer(prefix, key string) {
	r.serverMutex.Lock()
	defer r.serverMutex.Unlock()
	if _, ok := r.servers[prefix]; ok {
		delete(r.servers[prefix], key)
	}
}

func (r *Discovery) GetServers(prefix string) map[string]any {
	r.serverMutex.Lock()
	defer r.serverMutex.Unlock()
	return r.servers[prefix]
}
