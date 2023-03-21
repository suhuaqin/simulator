package node_client

import (
	"context"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"testing"
	"unknownName/proto/node"
)

func TestNewSenderClient(t *testing.T) {
	service := micro.NewService(
		micro.Name("client"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs("127.0.0.1:2379"),
		)),
	)

	c := NewNodeClient(service)
	_, err := c.SetInterval(context.Background(), &node.SetIntervalRequest{
		IntervalMillisecond: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	service.Options().Client.Init()

	//_, err := c.Send(context.Background(), &sender.SendRequest{
	//	Message:    []byte("dd"),
	//	ReceiverId: "receiver-914eb0bb-b91a-4f3b-841e-9b93cef95d2e",
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
}
