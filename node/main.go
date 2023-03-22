package main

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"simulator/client/node_client"
	"simulator/client/transfer_client"
	"simulator/frame/logWrapper"
	"simulator/node/handler"
	"simulator/node/repo"
	pb "simulator/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {
	etcdEndpoint := "127.0.0.1:2379"

	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(node_client.NodeServiceName),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(etcdEndpoint),
		)),
		micro.WrapHandler(logWrapper.LogWrapper),
	)

	etcdCli, err := repo.NewEtcdCli(etcdEndpoint)
	if err != nil {
		panic(err)
	}
	senderRepo := repo.NewNodeRepo(etcdCli)

	// Register handler
	nodeService := handler.NewNodeService(srv, senderRepo, transfer_client.NewTransferClient(srv))
	defer nodeService.Stop()
	if err := pb.RegisterNodeHandler(srv.Server(), nodeService); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
