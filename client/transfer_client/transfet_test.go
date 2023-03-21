package transfer_client

import (
	"context"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"testing"
	"unknownName/frame/debugfmt"
	"unknownName/proto/admin"
)

func TestNewTransferClient(t *testing.T) {
	service := micro.NewService(
		micro.Name("client"),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs("127.0.0.1:2379"),
		)),
	)

	c := admin.NewAdminService("admin", service.Client())
	//client := NewTransferClient(service)
	//_, err := client.Transfer(context.Background(), &transfer.TransferRequest{
	//	MsgId:      0,
	//	Message:    nil,
	//	ReceiverId: "",
	//	SenderId:   "",
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}

	//_, err := client.SetDiscard(context.Background(), &transfer.SetDiscardRequest{
	//	Remainder: 2,
	//	DiscardLe: 1,
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	repson, err := c.GetServices(context.Background(), &admin.GetServicesRequest{})
	if err != nil {
		t.Fatal(err)
	}
	debugfmt.JsonMarshalIndent(repson, "")
}
