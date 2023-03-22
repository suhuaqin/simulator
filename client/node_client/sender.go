package node_client

import (
	"go-micro.dev/v4"
	pb "simulator/proto"
)

const NodeServiceName = "node"

func NewNodeClient(srv micro.Service) pb.NodeService {
	s := pb.NewNodeService(NodeServiceName, srv.Client())
	return s
}
