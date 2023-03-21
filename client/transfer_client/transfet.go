package transfer_client

import (
	"go-micro.dev/v4"
	pb "simulator/proto"
)

const TransferServiceName = "transfer"

func NewTransferClient(srv micro.Service) pb.TransferService {
	s := pb.NewTransferService(TransferServiceName, srv.Client())
	return s
}
