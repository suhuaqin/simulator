package repo

import (
	"context"
	"simulator/frame/etcdKey"
)

type TransferRepository interface {
	GetNode(ctx context.Context) (etcdKey.NodeRegistryDataMap, error)
}
