package locate

import (
	"context"
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	objects = make(map[string]struct{})
	mu      sync.Mutex
)

func Locate(key string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/objects/%s", global.Cfg.Server.BaseDir, key))
	log.Println("Locate from local", key, err == nil)
	return err == nil
}

func StartLocate() {
	watchCh := global.Etcd.Watch(context.Background(), "/locates/request/", etcdClientv3.WithPrefix())
	for resp := range watchCh {
		for _, ev := range resp.Events {
			if ev.Type == etcdClientv3.EventTypePut {
				key := string(ev.Kv.Key)
				key = strings.Split(key, "/")[3]
				log.Println("Receive Locate request", key)
				if Locate(key) {
					values := strings.Split(string(ev.Kv.Value), ":")
					timestamp := values[1]
					k := fmt.Sprintf("/locates/result/%s/%s", key, timestamp)
					v := fmt.Sprintf("%s:%d", global.Cfg.Ip, global.Cfg.Port)
					log.Println("Locate result", k, v)
					_, err := global.Etcd.Put(context.Background(), k, v)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

func Load() {
	files, err := filepath.Glob(fmt.Sprintf("%s/objects/*", global.Cfg.Server.BaseDir))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		hash := filepath.Base(file)
		objects[hash] = struct{}{}
	}
}

func Add(key string) {
	mu.Lock()
	defer mu.Unlock()
	objects[key] = struct{}{}
}

func Del(hash string) {
	mu.Lock()
	defer mu.Unlock()
	delete(objects, hash)
}
