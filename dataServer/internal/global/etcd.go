package global

import (
	etcdClientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var Etcd *etcdClientv3.Client

func initEtcd() {
	var err error
	Etcd, err = etcdClientv3.New(etcdClientv3.Config{
		Endpoints:   Cfg.Etcd.Endpoints,
		DialTimeout: time.Duration(Cfg.Etcd.DialTimeout) * time.Second,
	})
	if err != nil {
		panic(err)
	}
}
