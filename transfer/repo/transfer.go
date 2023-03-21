package repo

import (
	"context"
	pb "simulator/proto"
)

type TransferRepository interface {
	GetReceiverRegistry(ctx context.Context) (*pb.ReceiverRegistry, error)
}
