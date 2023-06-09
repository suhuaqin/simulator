package repo

import (
	"context"
	"encoding/json"
	"go-micro.dev/v4/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"simulator/frame/etcdKey"
	"time"
)

type TransferRepo struct {
	cli              *clientv3.Client
	nodeRegistryData etcdKey.NodeRegistryDataMap
}

func NewTransferRepo(etcdCli *clientv3.Client) *TransferRepo {
	result := &TransferRepo{cli: etcdCli}
	result.nodeRegistryData, _ = result.getNode(context.Background())
	result.watchNode(context.Background())
	return result
}

func (s *TransferRepo) GetNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error) {
	return s.nodeRegistryData, nil
}

func (s *TransferRepo) watchNode(ctx context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[panic]", err)
			}
		}()

	reset:
		duration := 30 * time.Second
		watchCh := s.cli.Watch(ctx, etcdKey.NodeRegistryKeyPrefix, clientv3.WithPrefix())
		ticker := time.NewTicker(duration)
		for {
			select {
			case resp := <-watchCh:
				if resp.Err() != nil {
					logger.Error(resp.Err())
					goto reset
				}
				logger.Debug("[watch] nodeRegistryData")
				ticker.Reset(duration)
			case <-ticker.C:
			}
			receiverData, err := s.getNode(ctx)
			if err != nil {
				logger.Error(err)
			}
			s.nodeRegistryData = receiverData
		}
	}()
}

func (s *TransferRepo) getNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error) {
	result := make(etcdKey.NodeRegistryDataMap)
	resp, err := s.cli.Get(ctx, etcdKey.NodeRegistryKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return result, err
	}

	for _, kv := range resp.Kvs {
		resultElem := make(etcdKey.NodeRegistryDataMap)
		err := json.Unmarshal(kv.Value, &resultElem)
		if err != nil {
			return result, err
		}
		for k, v := range resultElem {
			result[k] = v
		}
	}
	return result, nil
}
