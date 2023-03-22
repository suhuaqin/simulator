package handler

import "sync"

// TODO：切分锁
type msgIDHelp struct {
	l sync.Mutex
	// 维护两个端点的消息ID,key为远端的node id
	msgID map[string]uint32
}

func NewMsgIDHelp() *msgIDHelp {
	return &msgIDHelp{
		msgID: make(map[string]uint32),
	}
}

func (m *msgIDHelp) NewMsgID(receiverID string) uint32 {
	m.l.Lock()
	defer m.l.Unlock()
	m.msgID[receiverID]++
	return m.msgID[receiverID]
}
