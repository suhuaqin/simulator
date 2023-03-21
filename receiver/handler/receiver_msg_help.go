package handler

import (
	"sync"
	"time"
)

type receiverMsgHelp struct {
	l sync.Mutex
	// key: nodeID
	lack map[string]*lack
}

type lack struct {
	maxID    uint32
	totalLen uint64
	lastTime time.Time
	// key: 缺失的消息ID
	lack map[uint32]struct{}
}

func newReceiverMsgHelp() *receiverMsgHelp {
	return &receiverMsgHelp{
		lack: make(map[string]*lack),
	}
}

func (h *receiverMsgHelp) SetMsgID(nodeID string, msgID uint32, msgLen int) {
	h.l.Lock()
	defer h.l.Unlock()
	if _, ok := h.lack[nodeID]; !ok {
		h.lack[nodeID] = &lack{
			maxID:    0,
			totalLen: 0,
			lack:     make(map[uint32]struct{}, 1),
		}
	}

	lack := h.lack[nodeID]
	lack.totalLen += uint64(msgLen)
	lack.lastTime = time.Now()

	// 缺失的消息
	for i := lack.maxID + 1; i < msgID; i++ {
		lack.lack[i] = struct{}{}
	}

	if msgID > lack.maxID {
		lack.maxID = msgID
	} else {
		if _, ok := lack.lack[msgID]; !ok {
			// 重复消息
			lack.totalLen -= uint64(msgLen)
		}
		delete(lack.lack, msgID)
	}
}
