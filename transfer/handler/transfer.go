package handler

import (
	"context"
	"errors"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"simulator/client/fitter"
	"simulator/client/node_client"
	pb "simulator/proto"
	"simulator/transfer/repo"
)

type TransferService struct {
	nodeCli      pb.NodeService
	transferRepo repo.TransferRepository
	transferHelp *transferHelp
}

func NewTransferService(srv micro.Service, transferRepo repo.TransferRepository) *TransferService {
	transferService := &TransferService{
		nodeCli:      node_client.NewNodeClient(srv),
		transferRepo: transferRepo,
		transferHelp: newTransferHelp(),
	}
	return transferService
}

func (m *TransferService) Transfer(ctx context.Context, request *pb.TransferRequest, response *pb.TransferResponse) error {
	if m.transferHelp.isDiscard() {
		logger.Infof("msg discard: %+v", *request)
		return nil
	}

	receiverDataMap, err := m.transferRepo.GetNode(ctx)
	if err != nil {
		return err
	}
	nodeData, exist := receiverDataMap.GetNode(request.ReceiverId)
	if !exist {
		return errors.New("can not found receiver")
	}

	_, err = m.nodeCli.Recv(ctx, &pb.RecvRequest{
		MsgId:      request.MsgId,
		Message:    request.Message,
		SenderId:   request.SenderId,
		ReceiverId: request.ReceiverId,
	}, fitter.ServiceID(nodeData.ServiceID))
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (m *TransferService) Broadcast(ctx context.Context, request *pb.BroadcastRequest, response *pb.BroadcastResponse) error {
	receiverDataMap, err := m.transferRepo.GetNode(ctx)
	if err != nil {
		return err
	}
	recvBroadcastReq := &pb.RecvBroadcastRequest{
		SenderId: request.SenderId,
		Message:  request.Message,
	}
	alreadyDo := make(map[string]struct{})
	for _, v := range receiverDataMap {
		if _, ok := alreadyDo[v.ServiceID]; ok {
			continue
		}
		alreadyDo[v.ServiceID] = struct{}{}
		go func(serviceID string) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("[panic]", err)
				}
			}()
			_, err := m.nodeCli.RecvBroadcast(ctx, recvBroadcastReq, fitter.ServiceID(serviceID))
			if err != nil {
				logger.Error()
			}
		}(v.ServiceID)
	}
	return nil
}

func (m *TransferService) SetDiscard(ctx context.Context, request *pb.SetDiscardRequest, response *pb.SetDiscardResponse) error {
	return m.transferHelp.setConfig(request.Remainder, request.DiscardLe)
}
