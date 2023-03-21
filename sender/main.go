package main

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"simulator/client/node_client"
	"simulator/client/transfer_client"
	pb "simulator/proto"
	"simulator/sender/handler"
	"simulator/sender/repo"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {
	etcdEndpoint := "127.0.0.1:2379"

	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(node_client.SenderServiceName),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(etcdEndpoint),
		)),
	)

	etcdCli, err := repo.NewEtcdCli(etcdEndpoint)
	if err != nil {
		panic(err)
	}
	senderRepo := repo.NewSenderRepo(etcdCli)

	// Register handler
	senderService := handler.NewSenderService(srv, senderRepo, transfer_client.NewTransferClient(srv))
	defer senderService.Stop()
	if err := pb.RegisterSenderHandler(srv.Server(), senderService); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
