package repo

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	pb "simulator/proto"
	"testing"
	"time"
)

var (
	cli        *clientv3.Client
	senderRepo *SenderRepo
	ctx        = context.Background()
)

func TestMain(m *testing.M) {
	var err error
	cli, err = NewEtcdCli("127.0.0.1:2379")
	if err != nil {
		panic(err)
	}
	senderRepo = NewTransferRepo(cli)
	time.Sleep(5 * time.Second)
	m.Run()
}

func TestNodeModeRepo_SetSender(t *testing.T) {
	err := senderRepo.SetSender(ctx, "node1", 1, &pb.SenderRegistry{
		SenderRegistry: map[string]*pb.SenderMode{
			"senderID": {
				NodeId:              "node1",
				IntervalMillisecond: 5,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
