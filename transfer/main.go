package main

import (
	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"simulator/client/transfer_client"
	"simulator/frame/logWrapper"
	pb "simulator/proto"
	"simulator/transfer/handler"
	"simulator/transfer/repo"
)

func main() {
	logger.Init(logger.WithLevel(logger.DebugLevel))
	etcdEndpoint := "127.0.0.1:2379"

	// create a new service
	service := micro.NewService(
		micro.Name(transfer_client.TransferServiceName),
		micro.Registry(etcd.NewRegistry(
			registry.Addrs(etcdEndpoint),
		)),
		micro.WrapHandler(logWrapper.LogWrapper),
	)

	// initialise flags
	service.Init()

	etcdCli, err := repo.NewEtcdCli(etcdEndpoint)
	if err != nil {
		panic(err)
	}
	transferRepo := repo.NewTransferRepo(etcdCli)

	// Register handler
	if err := pb.RegisterTransferHandler(service.Server(), handler.NewTransferService(service, transferRepo)); err != nil {
		logger.Fatal(err)
	}

	// start the service
	logger.Error(service.Run())
}
