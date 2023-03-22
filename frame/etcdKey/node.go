package etcdKey

import (
	"fmt"
	"math/rand"
)

const (
	CustomKeyPrefix       = "/micro/custom"
	NodeRegistryKeyPrefix = "/micro/custom/node"
)

func NodeRegistryKey(nodeID string) string {
	return fmt.Sprintf("%s/%s", NodeRegistryKeyPrefix, nodeID)
}

type NodeRegistryDataMap map[string]NodeRegistryData

type NodeRegistryData struct {
	ID        string
	ServiceID string
}

// nodeID == "" 时随机获取一个
func (n NodeRegistryDataMap) GetNode(id string) (*NodeRegistryData, bool) {
	if id != "" {
		result, exist := n[id]
		return &result, exist
	}

	l := len(n)
	if l == 0 {
		return nil, false
	}
	i := rand.Intn(l)
	for _, v := range n {
		if i == 0 {
			return &v, true
		}
		i--
	}
	return nil, false
}
