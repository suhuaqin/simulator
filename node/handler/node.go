package handler

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"simulator/node/repo"
	pb "simulator/proto"
)

type NodeService struct {
	nodeHelp     *NodeHelp
	nodeRepo     repo.NodeRepository
	transportCli pb.TransferService
}

func NewNodeService(srv micro.Service, nodeRepo repo.NodeRepository, transportCli pb.TransferService) *NodeService {
	result := &NodeService{
		nodeHelp:     NewNodeHelp(fmt.Sprintf("%s-%s", srv.Server().Options().Name, srv.Server().Options().Id), nodeRepo),
		nodeRepo:     nodeRepo,
		transportCli: transportCli,
	}
	return result
}

func (s *NodeService) Send(ctx context.Context, request *pb.SendRequest, response *pb.SendResponse) error {
	senderNode, exist := s.nodeHelp.GetNode("")
	if !exist {
		return errors.New("can not found senderNode")
	}
	receiverRegistryData, err := s.nodeHelp.GetRegistry(request.ReceiverId)
	if err != nil {
		return err
	}

	_, err = s.transportCli.Transfer(ctx, &pb.TransferRequest{
		MsgId:      senderNode.NewSenderMsgID(receiverRegistryData.ID),
		Message:    request.Message,
		ReceiverId: receiverRegistryData.ID,
		SenderId:   senderNode.id,
	})
	return err
}

func (s *NodeService) SendBroadcast(ctx context.Context, request *pb.RecvBroadcastRequest, response *pb.RecvBroadcastResponse) error {
	node, exist := s.nodeHelp.GetNode(request.SenderId)
	if !exist {
		return errors.New("sender can not found")
	}
	_, err := s.transportCli.Broadcast(ctx, &pb.BroadcastRequest{
		SenderId: node.id,
		Message:  request.Message,
	})
	return err
}

func (s *NodeService) Recv(ctx context.Context, request *pb.RecvRequest, response *pb.RecvResponse) error {
	node, exist := s.nodeHelp.GetNode(request.ReceiverId)
	if !exist {
		return errors.New("can not found senderNode")
	}
	node.receiverMsgHelp.SetMsgID(request.SenderId, request.MsgId, len(request.Message))
	return nil
}

func (s *NodeService) RecvBroadcast(ctx context.Context, request *pb.RecvBroadcastRequest, response *pb.RecvBroadcastResponse) error {
	return nil
}

func (s *NodeService) SetNodeCount(ctx context.Context, request *pb.SetNodeCountRequest, response *pb.SetNodeCountResponse) error {
	s.nodeHelp.SetNodeCount(request.Num, s)
	return nil
}

func (s *NodeService) SetInterval(ctx context.Context, request *pb.SetIntervalRequest, response *pb.SetIntervalResponse) error {
	s.nodeHelp.SetSenderInterval(request.SenderId, request.IntervalMillisecond)
	return nil
}

func (s *NodeService) Stop() {
	s.nodeHelp.Stop()
}
