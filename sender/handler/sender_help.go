package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"math/rand"
	pb "simulator/proto"
	"simulator/sender/repo"
	"sync"
	"time"
)

type SenderHelp struct {
	// 目前 lock sender
	sync.RWMutex
	sender map[string]*senderSimulator

	nodeID     string
	senderRepo repo.SenderRepository
	registryCh chan struct{}
}

func NewSenderHelp(nodeID string, senderRepo repo.SenderRepository) *SenderHelp {
	result := &SenderHelp{
		sender:     make(map[string]*senderSimulator, 0),
		senderRepo: senderRepo,
		registryCh: make(chan struct{}, 100),
		nodeID:     nodeID,
	}
	go func() {
		result.registry()
	}()
	return result
}

func (s *SenderHelp) SetSenderNum(num uint32, senderHandler pb.SenderHandler) {
	s.Lock()
	defer func() { s.registryCh <- struct{}{} }()
	defer s.Unlock()

	// 数量减少
	if len(s.sender) > int(num) {
		for k, v := range s.sender {
			if int(num) >= len(s.sender) {
				return
			}
			v.stop()
			delete(s.sender, k)
		}
	}

	// 数量增加
	for i := len(s.sender); i < int(num); i = len(s.sender) {
		sender := NewSenderSimulator(senderHandler)
		s.sender[sender.id] = sender
	}
	return
}

// senderID == "" 时随机挑选
func (s *SenderHelp) GetSender(senderID string) (*senderSimulator, bool) {
	s.RLock()
	defer s.RUnlock()

	if senderID != "" {
		result, ok := s.sender[senderID]
		return result, ok
	}

	l := len(s.sender)
	if l == 0 {
		return nil, false
	}
	index := rand.Int() % l
	i := 0
	for _, v := range s.sender {
		if i == index {
			return v, true
		}
		i++
	}
	return nil, false
}

func (s *SenderHelp) SetSenderInterval(senderID string, intervalMillisecond uint32) error {
	s.Lock()
	defer s.Unlock()

	if s.sender[senderID] != nil {
		select {
		case s.sender[senderID].intervalMillisecondCh <- intervalMillisecond:
			s.sender[senderID].intervalMillisecond = intervalMillisecond
			s.registryCh <- struct{}{}
		case <-time.After(200 * time.Millisecond):
			return errors.New("time out")
		}
	}
	return nil
}

func (s *SenderHelp) registry() {
	registryData := &pb.SenderRegistry{
		SenderRegistry: make(map[string]*pb.SenderMode),
	}
	registryTTL := int64(30)
	for {
		select {
		case <-time.After(time.Duration(registryTTL) * time.Second):
		case <-s.registryCh:
		readAllCh:
			for {
				select {
				case <-s.registryCh:
				default:
					break readAllCh
				}
			}
		}

		for k, _ := range registryData.SenderRegistry {
			delete(registryData.SenderRegistry, k)
		}
		s.RLock()
		for k, v := range s.sender {
			registryData.SenderRegistry[k] = &pb.SenderMode{
				NodeId:              s.nodeID,
				IntervalMillisecond: v.intervalMillisecond,
			}
		}
		s.RUnlock()
		for i := 0; i < 3; i++ {
			err := s.senderRepo.SetSender(context.Background(), s.nodeID, 3*registryTTL, registryData)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			break
		}
	}
}

func (s *SenderHelp) Stop() {
	s.senderRepo.DelSender(context.Background(), s.nodeID)
}

type senderSimulator struct {
	senderService         pb.SenderHandler
	baseCtx               context.Context
	cancel                context.CancelFunc
	id                    string
	intervalMillisecondCh chan uint32
	intervalMillisecond   uint32
	msgIDHelp             *msgIDHelp
}

func NewSenderSimulator(senderHandler pb.SenderHandler) *senderSimulator {
	result := &senderSimulator{
		senderService:         senderHandler,
		cancel:                nil,
		baseCtx:               nil,
		id:                    uuid.New().String(),
		intervalMillisecondCh: make(chan uint32, 1),
		msgIDHelp:             NewMsgIDHelp(),
	}
	result.baseCtx, result.cancel = context.WithCancel(context.Background())

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		result.demonsAutoSend()
	}()

	return result
}

func (s *senderSimulator) getMsgID(receiverID string) uint32 {
	return s.msgIDHelp.AddMsgID(receiverID)
}

func (s *senderSimulator) stop() {
	s.cancel()
}

func (s *senderSimulator) demonsAutoSend() {
	autoMsg := []byte("this is auto send message")
	ticker := time.NewTicker(time.Duration(<-s.intervalMillisecondCh) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := s.senderService.Send(context.Background(), &pb.SendRequest{
				Message: autoMsg,
			}, &pb.SendResponse{})
			if err != nil {
				logger.Error(err)
			}
		case interval := <-s.intervalMillisecondCh:
			ticker.Stop()
			if interval != 0 {
				ticker.Reset(time.Duration(interval) * time.Millisecond)
			}
		case <-s.baseCtx.Done():
			return
		}
	}
}
