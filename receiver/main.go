package main

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4/registry"
	"simulator/client/node_client"
	pb "simulator/proto"
	"simulator/receiver/handler"
	"simulator/receiver/repo"

	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
)

func main() {
	etcdEndpoint := "127.0.0.1:2379"

	// Create service
	srv := micro.NewService()
	srv.Init(
		micro.Name(node_client.ReceiverServiceName),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(etcdEndpoint),
		)),
	)

	etcdCli, err := repo.NewEtcdCli(etcdEndpoint)
	if err != nil {
		panic(err)
	}
	receiverRepo := repo.NewReceiverRepo(etcdCli)

	// Register handler
	receiverService := handler.NewReceiverService(srv, receiverRepo)
	defer receiverService.Stop()
	if err := pb.RegisterReceiverHandler(srv.Server(), receiverService); err != nil {
		logger.Fatal(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
