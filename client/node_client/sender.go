package node_client

import (
	"go-micro.dev/v4"
	pb "simulator/proto"
)

const SenderServiceName = "sender"
const ReceiverServiceName = "receiver"

func NewSenderClient(srv micro.Service) pb.SenderService {
	s := pb.NewSenderService(SenderServiceName, srv.Client())
	return s
}

func NewReceiverClient(srv micro.Service) pb.ReceiverService {
	s := pb.NewReceiverService(ReceiverServiceName, srv.Client())
	return s
}
