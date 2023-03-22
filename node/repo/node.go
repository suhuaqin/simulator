package repo

import (
	"context"
	"simulator/frame/etcdKey"
)

type NodeRepository interface {
	SetNode(ctx context.Context, nodeID string, ttlSecond int64, registry etcdKey.NodeRegistryDataMap) error
	GetNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error)
	DelNode(ctx context.Context, nodeID string) error
}
