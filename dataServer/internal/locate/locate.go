package locate

import (
	"context"
	"fmt"
	global2 "github.com/qiaofufu/tinyoss_kernal/dataServer/internal/global"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"strings"
	"time"
)

func Locate(key string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/objects/%s", global2.Cfg.Server.BaseDir, key))
	log.Println("Locate from local", key, err == nil)
	return err == nil
}

func LocateFromAllServer(key string) (string, error) {
	// publish locate information
	k := fmt.Sprintf("/locates/request/%s", key)
	timestamp := time.Now().UnixMilli()
	v := fmt.Sprintf("timestam:%d", timestamp)
	lease, err := global2.Etcd.Grant(context.Background(), global2.Cfg.Server.LocateTimeout)
	if err != nil {
		return "", err
	}
	log.Println("Locate request", k, v)
	_, err = global2.Etcd.Put(context.Background(), k, v, etcdClientv3.WithLease(lease.ID))
	if err != nil {
		return "", err
	}
	// wait for locate
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global2.Cfg.Server.LocateTimeout)*time.Second)
	defer cancel()
	resultKey := fmt.Sprintf("/locates/result/%s/%d", key, timestamp)
	resultCh := global2.Etcd.Watch(ctx, resultKey, etcdClientv3.WithPrefix())
	for wresp := range resultCh {
		for _, ev := range wresp.Events {
			if ev.Type == etcdClientv3.EventTypePut {
				return string(ev.Kv.Value), nil
			}
		}
	}
	return "", nil
}

func StartLocate() {
	watchCh := global2.Etcd.Watch(context.Background(), "/locates/request/", etcdClientv3.WithPrefix())
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
					v := fmt.Sprintf("%s:%d", global2.Cfg.Ip, global2.Cfg.Port)
					log.Println("Locate result", k, v)
					_, err := global2.Etcd.Put(context.Background(), k, v)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
