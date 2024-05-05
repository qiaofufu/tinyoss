package meta

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type Metadata struct {
	Name    string
	Hash    string
	Size    int64
	Version int
}

var etcdClient *clientv3.Client

func InitMeta(etcd *clientv3.Client) {
	etcdClient = etcd
}

func SearchAllVersions(object string) ([]*Metadata, error) {
	key := fmt.Sprintf("/meta/objects/%s", object)
	resp, err := etcdClient.Get(context.Background(), key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}
	var metas []*Metadata
	for _, kv := range resp.Kvs {
		meta := &Metadata{}
		err = json.Unmarshal(kv.Value, meta)
		if err != nil {
			return nil, err
		}
		metas = append(metas, meta)
	}
	return metas, nil
}

func GetLatestVersion(object string) (*Metadata, error) {
	metas, err := SearchAllVersions(object)
	if err != nil {
		return nil, err
	}
	if len(metas) == 0 {
		log.Printf("object version %s not found", object)
		return nil, nil
	}
	log.Printf("latest version: %v", metas[0])
	return metas[0], nil
}

func GetMetadata(object string, version int) (*Metadata, error) {
	if version == 0 {
		return GetLatestVersion(object)
	}
	key := fmt.Sprintf("/meta/objects/%s/%d", object, version)
	resp, err := etcdClient.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	meta := &Metadata{}
	if resp.Count == 0 {
		return nil, fmt.Errorf("object %s version %d not found", object, version)
	}
	err = json.Unmarshal(resp.Kvs[0].Value, meta)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func Del(name string) error {
	latest, err := GetLatestVersion(name)
	if err != nil {
		return err
	}
	if latest == nil {
		return nil
	}
	latest.Version = latest.Version + 1
	latest.Hash = ""
	latest.Size = 0
	return Put(name, latest)
}

func Put(name string, value *Metadata) error {
	key := fmt.Sprintf("/meta/objects/%s/%d", name, value.Version)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = etcdClient.Put(context.Background(), key, string(data))
	if err != nil {
		return err
	}
	return nil
}

func Add(name string, hash string, size int64) error {
	latest, err := GetLatestVersion(name)
	if err != nil {
		return err
	}
	version := 0
	if latest != nil {
		version = latest.Version + 1
	}
	meta := &Metadata{
		Name:    name,
		Hash:    hash,
		Size:    size,
		Version: version,
	}
	return Put(name, meta)
}
