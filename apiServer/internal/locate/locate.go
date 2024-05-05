package locate

import (
	"context"
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func LocateFromAllServer(key string) (string, error) {
	// publish locate information
	k := fmt.Sprintf("/locates/request/%s", key)
	timestamp := time.Now().UnixMilli()
	v := fmt.Sprintf("timestamp:%d", timestamp)
	lease, err := global.Etcd.Grant(context.Background(), global.Cfg.Server.LocateTimeout)
	if err != nil {
		return "", err
	}
	log.Println("Locate request", k, v)

	// wait for locate
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Cfg.Server.LocateTimeout)*time.Second)
	defer cancel()
	resultKey := fmt.Sprintf("/locates/result/%s/%d", key, timestamp)
	resultCh := global.Etcd.Watch(ctx, resultKey, etcdClientv3.WithPrefix())
	_, err = global.Etcd.Put(context.Background(), k, v, etcdClientv3.WithLease(lease.ID))
	if err != nil {
		return "", err
	}
	for wresp := range resultCh {
		for _, ev := range wresp.Events {
			if ev.Type == etcdClientv3.EventTypePut {
				return string(ev.Kv.Value), nil
			}
		}
	}
	return "", nil
}

func Exist(key string) (bool, error) {
	server, err := LocateFromAllServer(key)
	if err != nil {
		return false, err
	}
	return server != "", nil
}
