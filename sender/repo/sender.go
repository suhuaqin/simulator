package repo

import (
	"context"
	pb "simulator/proto"
)

type SenderRepository interface {
	SetSender(ctx context.Context, nodeID string, ttlSecond int64, registry *pb.SenderRegistry) error
	GetReceiverRegistry(ctx context.Context) (*pb.ReceiverRegistry, error)
	DelSender(ctx context.Context, nodeID string) error
}
