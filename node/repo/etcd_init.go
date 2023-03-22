package repo

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func NewEtcdCli(endpoint string) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 2 * time.Second,
	})
	return cli, err
}
