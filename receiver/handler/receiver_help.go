package handler

import (
	"context"
	"github.com/google/uuid"
	pb "simulator/proto"
	"simulator/receiver/repo"
	"sync"
	"time"
)

type receiverHelp struct {
	// 目前 lock receiver
	sync.RWMutex
	receiver map[string]*receiverSimulator

	nodeID       string
	receiverRepo repo.ReceiverRepository
	registryCh   chan struct{}
}

func NewReceiverHelp(nodeID string, receiverRepo repo.ReceiverRepository) *receiverHelp {
	result := &receiverHelp{
		receiver:     make(map[string]*receiverSimulator, 0),
		nodeID:       nodeID,
		receiverRepo: receiverRepo,
		registryCh:   make(chan struct{}, 100),
	}
	go func() {
		result.registry()
	}()
	return result
}

func (r *receiverHelp) SetMsgID(senderID, receiverID string, msgID uint32, msgLen int) {
	r.Lock()
	defer r.Unlock()

	r.receiver[receiverID].msgHelp.SetMsgID(senderID, msgID, msgLen)
}

func (r *receiverHelp) SetNum(num uint32) {
	r.Lock()
	defer func() { r.registryCh <- struct{}{} }()
	defer r.Unlock()

	// 数量减少
	if len(r.receiver) > int(num) {
		for k, _ := range r.receiver {
			if int(num) >= len(r.receiver) {
				return
			}
			delete(r.receiver, k)
		}
	}

	// 数量增加
	for i := len(r.receiver); i < int(num); i = len(r.receiver) {
		receiver := NewReceiverSimulator()
		r.receiver[receiver.id] = receiver
	}
	return
}

func (r *receiverHelp) registry() {
	registryData := &pb.ReceiverRegistry{
		ReceiverRegistry: make(map[string]string),
	}
	registryTTL := int64(30)
	for {
		select {
		case <-time.After(time.Duration(registryTTL) * time.Second):
		case <-r.registryCh:
		readAllCh:
			for {
				select {
				case <-r.registryCh:
				default:
					break readAllCh
				}
			}
		}

		for k, _ := range registryData.ReceiverRegistry {
			delete(registryData.ReceiverRegistry, k)
		}
		r.RLock()
		for k, _ := range r.receiver {
			registryData.ReceiverRegistry[k] = r.nodeID
		}
		r.RUnlock()
		for i := 0; i < 3; i++ {
			err := r.receiverRepo.SetReceiver(context.Background(), r.nodeID, 3*registryTTL, registryData)
			if err != nil {
				time.Sleep(100 * time.Millisecond)
			}
			break
		}
	}
}

func (r *receiverHelp) Stop() {
	r.receiverRepo.DelSender(context.Background(), r.nodeID)
}

type receiverSimulator struct {
	id      string
	msgHelp *receiverMsgHelp
}

func NewReceiverSimulator() *receiverSimulator {
	result := &receiverSimulator{
		id:      uuid.New().String(),
		msgHelp: newReceiverMsgHelp(),
	}
	return result
}
