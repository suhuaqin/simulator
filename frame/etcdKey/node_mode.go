package etcdKey

import (
	"encoding/json"
	"fmt"
	pb "simulator/proto"
)

const (
	customKeyPrefix     = "micro/custom"
	SenderRegistryKey   = "micro/custom/sender/"
	ReceiverRegistryKey = "micro/custom/receiver/"
)

func SenderRegistryNodeKey(nodeID string) string {
	return fmt.Sprintf("%s%s", SenderRegistryKey, nodeID)
}

func ReceiverRegistryNodeKey(nodeID string) string {
	return fmt.Sprintf("%s%s", ReceiverRegistryKey, nodeID)
}

func SenderMarshal(nodeMode *pb.SenderRegistry) ([]byte, error) {
	return json.Marshal(nodeMode)
}

func SenderUnmarshal(by []byte) (*pb.SenderRegistry, error) {
	result := &pb.SenderRegistry{
		SenderRegistry: make(map[string]*pb.SenderMode),
	}
	if len(by) == 0 {
		return result, nil
	}

	err := json.Unmarshal(by, &result)
	return result, err
}

func ReceiverMarshal(nodeMode *pb.ReceiverRegistry) ([]byte, error) {
	return json.Marshal(nodeMode)
}

func ReceiverUnmarshal(by []byte) (*pb.ReceiverRegistry, error) {
	result := &pb.ReceiverRegistry{
		ReceiverRegistry: make(map[string]string),
	}
	if len(by) == 0 {
		return result, nil
	}

	err := json.Unmarshal(by, &result)
	return result, err
}
