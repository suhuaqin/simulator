package repo

import (
	"context"
	pb "simulator/proto"
)

type ReceiverRepository interface {
	SetReceiver(ctx context.Context, nodeID string, ttlSecond int64, registry *pb.ReceiverRegistry) error
	DelSender(ctx context.Context, nodeID string) error
}
