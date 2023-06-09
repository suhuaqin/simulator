// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/transfer.proto

package simulator

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Transfer service

func NewTransferEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Transfer service

type TransferService interface {
	Transfer(ctx context.Context, in *TransferRequest, opts ...client.CallOption) (*TransferResponse, error)
	SetDiscard(ctx context.Context, in *SetDiscardRequest, opts ...client.CallOption) (*SetDiscardResponse, error)
	Broadcast(ctx context.Context, in *BroadcastRequest, opts ...client.CallOption) (*BroadcastResponse, error)
}

type transferService struct {
	c    client.Client
	name string
}

func NewTransferService(name string, c client.Client) TransferService {
	return &transferService{
		c:    c,
		name: name,
	}
}

func (c *transferService) Transfer(ctx context.Context, in *TransferRequest, opts ...client.CallOption) (*TransferResponse, error) {
	req := c.c.NewRequest(c.name, "Transfer.Transfer", in)
	out := new(TransferResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transferService) SetDiscard(ctx context.Context, in *SetDiscardRequest, opts ...client.CallOption) (*SetDiscardResponse, error) {
	req := c.c.NewRequest(c.name, "Transfer.SetDiscard", in)
	out := new(SetDiscardResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transferService) Broadcast(ctx context.Context, in *BroadcastRequest, opts ...client.CallOption) (*BroadcastResponse, error) {
	req := c.c.NewRequest(c.name, "Transfer.Broadcast", in)
	out := new(BroadcastResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Transfer service

type TransferHandler interface {
	Transfer(context.Context, *TransferRequest, *TransferResponse) error
	SetDiscard(context.Context, *SetDiscardRequest, *SetDiscardResponse) error
	Broadcast(context.Context, *BroadcastRequest, *BroadcastResponse) error
}

func RegisterTransferHandler(s server.Server, hdlr TransferHandler, opts ...server.HandlerOption) error {
	type transfer interface {
		Transfer(ctx context.Context, in *TransferRequest, out *TransferResponse) error
		SetDiscard(ctx context.Context, in *SetDiscardRequest, out *SetDiscardResponse) error
		Broadcast(ctx context.Context, in *BroadcastRequest, out *BroadcastResponse) error
	}
	type Transfer struct {
		transfer
	}
	h := &transferHandler{hdlr}
	return s.Handle(s.NewHandler(&Transfer{h}, opts...))
}

type transferHandler struct {
	TransferHandler
}

func (h *transferHandler) Transfer(ctx context.Context, in *TransferRequest, out *TransferResponse) error {
	return h.TransferHandler.Transfer(ctx, in, out)
}

func (h *transferHandler) SetDiscard(ctx context.Context, in *SetDiscardRequest, out *SetDiscardResponse) error {
	return h.TransferHandler.SetDiscard(ctx, in, out)
}

func (h *transferHandler) Broadcast(ctx context.Context, in *BroadcastRequest, out *BroadcastResponse) error {
	return h.TransferHandler.Broadcast(ctx, in, out)
}
