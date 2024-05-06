package locate

import (
	"context"
	"fmt"
	"github.com/qiaofufu/tinyoss_kernal/apiServer/internal/global"
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

func LocateFromAllServer(hash string) (map[int]string, error) {
	locateInfo := make(map[int]string)
	wg := sync.WaitGroup{}
	wg.Add(global.Cfg.RS.ShardAllNum)
	mu := sync.Mutex{}
	for i := 0; i < global.Cfg.RS.ShardAllNum; i++ {
		go func(i int) {
			defer wg.Done()
			k := fmt.Sprintf("/locates/request/%s.%d", hash, i)
			timestamp := time.Now().UnixMilli()
			v := fmt.Sprintf("timestamp:%d", timestamp)

			lease, err := global.Etcd.Grant(context.Background(), global.Cfg.Server.LocateTimeout)
			if err != nil {
				return
			}
			// wait for locate
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Cfg.Server.LocateTimeout)*time.Second)
			defer cancel()
			resultKey := fmt.Sprintf("/locates/result/%s/%d/%d", hash, i, timestamp)
			resultCh := global.Etcd.Watch(ctx, resultKey, etcdClientv3.WithPrefix())
			_, err = global.Etcd.Put(context.Background(), k, v, etcdClientv3.WithLease(lease.ID))
			if err != nil {
				return
			}
			for wresp := range resultCh {
				for _, ev := range wresp.Events {
					if ev.Type == etcdClientv3.EventTypePut {
						mu.Lock()
						locateInfo[i] = string(ev.Kv.Value)
						mu.Unlock()
					}
				}
			}
		}(i)
	}
	wg.Wait()
	return locateInfo, nil
}

func Exist(key string) (bool, error) {
	servers, err := LocateFromAllServer(key)
	if err != nil {
		return false, err
	}
	return len(servers) >= global.Cfg.RS.DataShard, nil
}
