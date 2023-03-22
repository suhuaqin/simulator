package repo

import (
	"context"
	"encoding/json"
	"go-micro.dev/v4/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"simulator/frame/etcdKey"
	"time"
)

type nodeRepo struct {
	cli              *clientv3.Client
	nodeRegistryData etcdKey.NodeRegistryDataMap
}

func NewNodeRepo(etcdCli *clientv3.Client) *nodeRepo {
	result := &nodeRepo{cli: etcdCli}
	result.nodeRegistryData, _ = result.getNode(context.Background())
	result.watchNode(context.Background())
	return result
}

func (s *nodeRepo) SetNode(ctx context.Context, nodeID string, ttlSecond int64, registry etcdKey.NodeRegistryDataMap) error {
	lease, err := s.cli.Grant(ctx, ttlSecond)
	if err != nil {
		return err
	}
	by, err := json.Marshal(registry)
	if err != nil {
		return err
	}
	_, err = s.cli.Put(ctx, etcdKey.NodeRegistryKey(nodeID), string(by), clientv3.WithLease(lease.ID))
	return err
}

func (s *nodeRepo) GetNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error) {
	return s.nodeRegistryData, nil
}

func (s *nodeRepo) DelNode(ctx context.Context, nodeID string) error {
	_, err := s.cli.Delete(ctx, etcdKey.NodeRegistryKey(nodeID))
	return err
}

func (s *nodeRepo) watchNode(ctx context.Context) {
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

func (s *nodeRepo) getNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error) {
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
