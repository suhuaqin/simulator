package repo

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"simulator/frame/etcdKey"
	pb "simulator/proto"
)

type ReceiverRepo struct {
	cli *clientv3.Client
}

func NewReceiverRepo(etcdCli *clientv3.Client) *ReceiverRepo {
	return &ReceiverRepo{cli: etcdCli}
}

func (s *ReceiverRepo) SetReceiver(ctx context.Context, nodeID string, ttlSecond int64, registry *pb.ReceiverRegistry) error {
	lease, err := s.cli.Grant(ctx, ttlSecond)
	if err != nil {
		return err
	}
	by, err := etcdKey.ReceiverMarshal(registry)
	if err != nil {
		return err
	}
	_, err = s.cli.Put(ctx, etcdKey.ReceiverRegistryNodeKey(nodeID), string(by), clientv3.WithLease(lease.ID))
	return err
}

func (s *ReceiverRepo) DelSender(ctx context.Context, nodeID string) error {
	_, err := s.cli.Delete(ctx, etcdKey.ReceiverRegistryNodeKey(nodeID))
	return err
}
