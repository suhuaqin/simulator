package fitter

import (
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/selector"
)

func NodeID(nodeID string) client.CallOption {
	filter := func(services []*registry.Service) []*registry.Service {
		var filtered []*registry.Service

		for _, service := range services {
			for _, node := range service.Nodes {
				if node.Id == nodeID {
					filtered = append(filtered, &registry.Service{
						Name:      service.Name,
						Version:   service.Version,
						Metadata:  service.Metadata,
						Endpoints: service.Endpoints,
						Nodes:     []*registry.Node{node},
					})
				}
			}
		}

		return filtered
	}

	return client.WithSelectOption(selector.WithFilter(filter))
}
