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
	receiverClient pb.ReceiverService
	transferRepo   repo.TransferRepository
	transferHelp   *transferHelp
}

func NewTransferService(srv micro.Service, transferRepo repo.TransferRepository) *TransferService {
	transferService := &TransferService{
		receiverClient: node_client.NewReceiverClient(srv),
		transferRepo:   transferRepo,
		transferHelp:   newTransferHelp(),
	}
	return transferService
}

func (m *TransferService) Transfer(ctx context.Context, request *pb.TransferRequest, response *pb.TransferResponse) error {
	if m.transferHelp.isDiscard() {
		logger.Infof("msg discard: %+v", *request)
		return nil
	}

	receiverData, err := m.transferRepo.GetReceiverRegistry(ctx)
	if err != nil {
		return err
	}
	nodeID := receiverData.GetReceiverRegistry()[request.GetReceiverId()]
	if nodeID == "" {
		return errors.New("can not found receiver")
	}

	_, err = m.receiverClient.Recv(ctx, &pb.RecvRequest{
		MsgId:      request.MsgId,
		Message:    request.Message,
		SenderId:   request.SenderId,
		ReceiverId: request.ReceiverId,
	}, fitter.NodeID(nodeID))
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (m *TransferService) SetDiscard(ctx context.Context, request *pb.SetDiscardRequest, response *pb.SetDiscardResponse) error {
	return m.transferHelp.setConfig(request.Remainder, request.DiscardLe)
}
