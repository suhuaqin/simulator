package repo

import (
	"context"
	"go-micro.dev/v4/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"simulator/frame/etcdKey"
	pb "simulator/proto"
	"time"
)

type TransferRepo struct {
	cli          *clientv3.Client
	receiverData *pb.ReceiverRegistry
}

func NewTransferRepo(etcdCli *clientv3.Client) *TransferRepo {
	result := &TransferRepo{cli: etcdCli}
	result.getReceiverRegistry(context.Background())
	result.watchReceiverRegistry(context.Background())
	return result
}

func (s *TransferRepo) GetReceiverRegistry(ctx context.Context) (*pb.ReceiverRegistry, error) {
	return s.receiverData, nil
}

func (s *TransferRepo) watchReceiverRegistry(ctx context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[panic]", err)
			}
		}()
	reset:
		duration := 30 * time.Second
		watchCh := s.cli.Watch(ctx, etcdKey.ReceiverRegistryKey, clientv3.WithPrefix())
		ticker := time.NewTicker(duration)
		for {
			select {
			case resp := <-watchCh:
				if resp.Err() != nil {
					logger.Error(resp.Err())
					goto reset
				}
				logger.Debug("[watch] receiverData")
				ticker.Reset(duration)
			case <-ticker.C:
			}
			receiverData, err := s.getReceiverRegistry(ctx)
			if err != nil {
				logger.Error(err)
			}
			s.receiverData = receiverData
		}
	}()
}

func (s *TransferRepo) getReceiverRegistry(ctx context.Context) (*pb.ReceiverRegistry, error) {
	result := &pb.ReceiverRegistry{
		ReceiverRegistry: make(map[string]string),
	}
	resp, err := s.cli.Get(ctx, etcdKey.ReceiverRegistryKey, clientv3.WithPrefix())
	if err != nil {
		return result, err
	}
	for _, kv := range resp.Kvs {
		nodeRegistry, err := etcdKey.ReceiverUnmarshal(kv.Value)
		if err != nil {
			return result, err
		}
		for k, v := range nodeRegistry.GetReceiverRegistry() {
			result.ReceiverRegistry[k] = v
		}
	}
	return result, nil
}
