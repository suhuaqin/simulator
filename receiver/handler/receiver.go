package handler

import (
	"context"
	"fmt"
	"go-micro.dev/v4"
	"simulator/frame/debugfmt"
	pb "simulator/proto"
	"simulator/receiver/repo"
)

type ReceiverService struct {
	receiverHelp *receiverHelp
}

func NewReceiverService(srv micro.Service, receiverRepo repo.ReceiverRepository) *ReceiverService {
	result := &ReceiverService{
		receiverHelp: NewReceiverHelp(fmt.Sprintf("%s-%s", srv.Server().Options().Name, srv.Server().Options().Id), receiverRepo),
	}
	return result
}

func (s *ReceiverService) Recv(ctx context.Context, request *pb.RecvRequest, response *pb.RecvResponse) error {
	debugfmt.JsonMarshalIndent(request, "")
	s.receiverHelp.SetMsgID(request.SenderId, request.ReceiverId, request.MsgId, len(request.Message))
	return nil
}

func (s *ReceiverService) SetReceiverNum(ctx context.Context, request *pb.SetReceiverNumRequest, response *pb.SetReceiverNumResponse) error {
	s.receiverHelp.SetNum(request.Num)
	return nil
}

func (s *ReceiverService) Stop() {
	s.receiverHelp.Stop()
}
