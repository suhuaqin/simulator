package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go-micro.dev/v4/codec"
	"go-micro.dev/v4/codec/json"
	"go-micro.dev/v4/logger"
	"math/rand"
	"simulator/frame/etcdKey"
	"simulator/node/repo"
	pb "simulator/proto"
	"sync"
	"time"
)

type NodeHelp struct {
	sync.RWMutex
	nodes map[string]*node

	nodeRepo   repo.NodeRepository
	serviceID  string
	registryCh chan struct{}
}

func NewNodeHelp(serviceID string, nodeRepo repo.NodeRepository) *NodeHelp {
	result := &NodeHelp{
		nodes:      make(map[string]*node),
		nodeRepo:   nodeRepo,
		serviceID:  serviceID,
		registryCh: make(chan struct{}, 10),
	}
	go func() {
		go func() {
			if err := recover(); err != nil {
				logger.Error("[panic]", err)
			}
		}()
		result.registry()
	}()
	return result
}

func (n *NodeHelp) SetNodeCount(num uint32, nodeHandler pb.NodeHandler) {
	n.Lock()
	defer func() { n.registryCh <- struct{}{} }()
	defer n.Unlock()

	// 数量减少
	if len(n.nodes) > int(num) {
		for k, v := range n.nodes {
			if int(num) >= len(n.nodes) {
				break
			}
			v.stop()
			delete(n.nodes, k)
		}
		return
	}

	// 数量增加
	for i := len(n.nodes); i < int(num); i = len(n.nodes) {
		node := newNode(json.Marshaler{}, nodeHandler)
		n.nodes[node.id] = node
	}
	return
}

// id == "" 时随机获取
func (n *NodeHelp) GetNode(id string) (*node, bool) {
	n.RLock()
	defer n.RUnlock()
	if id != "" {
		result, exist := n.nodes[id]
		return result, exist
	}

	l := len(n.nodes)
	if l == 0 {
		return nil, false
	}
	i := rand.Intn(l)
	for _, v := range n.nodes {
		if i == 0 {
			return v, true
		}
		i--
	}
	return nil, false
}

func (n *NodeHelp) GetRegistry(id string) (*etcdKey.NodeRegistryData, error) {
	nodeRegistryData, err := n.nodeRepo.GetNode(context.Background())
	if err != nil {
		return nil, err
	}

	nodeData, exist := nodeRegistryData.GetNode(id)
	if !exist {
		return nil, errors.New("can not found receiverNode")
	}
	return nodeData, nil
}

func (n *NodeHelp) registry() {
	registryDataMap := make(etcdKey.NodeRegistryDataMap)
	registryTTL := int64(30)
	for {
		select {
		case <-time.After(time.Duration(registryTTL) * time.Second):
		case <-n.registryCh:
		readAllCh:
			for {
				select {
				case <-n.registryCh:
				default:
					break readAllCh
				}
			}
		}

		for k, _ := range registryDataMap {
			delete(registryDataMap, k)
		}
		n.RLock()
		for k, _ := range n.nodes {
			registryDataMap[k] = etcdKey.NodeRegistryData{
				ID:        k,
				ServiceID: n.serviceID,
			}
		}
		n.RUnlock()
		for i := 0; i < 3; i++ {
			err := n.nodeRepo.SetNode(context.Background(), n.serviceID, 3*registryTTL, registryDataMap)
			if err != nil {
				logger.Error("[node_registry]", err)
				time.Sleep(100 * time.Millisecond)
			}
			break
		}
	}
}

func (n *NodeHelp) Stop() {
	n.nodeRepo.DelNode(context.Background(), n.serviceID)
}

func (n *NodeHelp) SetSenderInterval(nodeID string, intervalMillisecond uint32) error {
	node, exist := n.GetNode(nodeID)
	if !exist || node.id != nodeID {
		return nil
	}

	if n.nodes[nodeID] != nil {
		select {
		case n.nodes[nodeID].intervalMillisecondCh <- intervalMillisecond:
			n.nodes[nodeID].intervalMillisecond = intervalMillisecond
			n.registryCh <- struct{}{}
		case <-time.After(200 * time.Millisecond):
			return errors.New("time out")
		}
	}

	return nil
}

type node struct {
	id              string
	sendMsgHelp     *msgIDHelp
	marshaler       codec.Marshaler
	receiverMsgHelp *receiverMsgHelp

	intervalMillisecondCh chan uint32
	intervalMillisecond   uint32

	nodeHandler pb.NodeHandler

	baseCtx context.Context
	cancel  context.CancelFunc
}

func newNode(marshaler codec.Marshaler, nodeHandler pb.NodeHandler) *node {
	result := &node{
		id:                    uuid.New().String(),
		sendMsgHelp:           NewMsgIDHelp(),
		marshaler:             marshaler,
		receiverMsgHelp:       newReceiverMsgHelp(),
		intervalMillisecondCh: make(chan uint32, 1),
		intervalMillisecond:   0,
		nodeHandler:           nodeHandler,
	}
	result.baseCtx, result.cancel = context.WithCancel(context.Background())

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("[panic]", err)
			}
		}()
		result.demonsAutoSend()
	}()

	return result
}

func (n *node) NewSenderMsgID(receiverID string) uint32 {
	return n.sendMsgHelp.NewMsgID(receiverID)
}

func (n *node) Marshal(v any) ([]byte, error) {
	return n.marshaler.Marshal(v)
}

func (n *node) Unmarshal(by []byte, v any) error {
	return n.marshaler.Unmarshal(by, v)
}

func (n *node) demonsAutoSend() {
	autoMsg := []byte("this is auto send message")
	ticker := time.NewTicker(time.Duration(<-n.intervalMillisecondCh) * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			err := n.nodeHandler.Send(context.Background(), &pb.SendRequest{
				SenderId: n.id,
				Message:  autoMsg,
			}, &pb.SendResponse{})
			if err != nil {
				logger.Error(err)
			}
		case interval := <-n.intervalMillisecondCh:
			ticker.Stop()
			if interval != 0 {
				ticker.Reset(time.Duration(interval) * time.Millisecond)
			}
		case <-n.baseCtx.Done():
			return
		}
	}
}

func (n *node) stop() {
	n.cancel()
}
