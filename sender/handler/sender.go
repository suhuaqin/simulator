package handler

import (
	"context"
	"errors"
	"fmt"
	"go-micro.dev/v4"
	"math/rand"
	pb "simulator/proto"
	"simulator/sender/repo"
	"time"
)

type SenderService struct {
	senderHelp  *SenderHelp
	transferCli pb.TransferService

	senderRepo repo.SenderRepository
}

func NewSenderService(srv micro.Service, senderRepo repo.SenderRepository, transferCli pb.TransferService) *SenderService {
	rand.Seed(time.Now().Unix())
	result := &SenderService{
		senderHelp:  NewSenderHelp(fmt.Sprintf("%s-%s", srv.Server().Options().Name, srv.Server().Options().Id), senderRepo),
		transferCli: transferCli,
		senderRepo:  senderRepo,
	}
	return result
}

func (s *SenderService) SetSenderNum(ctx context.Context, request *pb.SetSenderNumRequest, response *pb.SetSenderNumResponse) error {
	s.senderHelp.SetSenderNum(request.Num, s)
	return nil
}

func (s *SenderService) Send(ctx context.Context, request *pb.SendRequest, response *pb.SendResponse) error {
	registryData, err := s.senderRepo.GetReceiverRegistry(ctx)
	if err != nil {
		return err
	}

	// 获取一个receiver
	receiverID := ""
	if request.ReceiverId != "" {
		receiverID = request.ReceiverId
	} else {
		l := len(registryData.GetReceiverRegistry())
		if l != 0 {
			index := rand.Int() % l
			i := 0
			for k, _ := range registryData.GetReceiverRegistry() {
				if i == index {
					receiverID = k
					break
				}
				i++
			}
		}
	}
	if receiverID == "" {
		return errors.New("can not found receiver")
	}

	// 获取一个sender
	sender, exist := s.senderHelp.GetSender("")
	if !exist {
		return errors.New("can not found sender")
	}

	_, err = s.transferCli.Transfer(ctx, &pb.TransferRequest{
		MsgId:      sender.getMsgID(receiverID),
		Message:    request.Message,
		ReceiverId: receiverID,
		SenderId:   sender.id,
	})

	return err
}

func (s *SenderService) SetInterval(ctx context.Context, request *pb.SetIntervalRequest, response *pb.SetIntervalResponse) error {
	return s.senderHelp.SetSenderInterval(request.GetSenderId(), request.GetIntervalMillisecond())
}

func (s *SenderService) Stop() {
	s.senderHelp.Stop()
}
