package logWrapper

import (
	"context"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
)

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		logger.Infof("[wrapper] server request: %v\nrequest: %+v", req.Endpoint(), req.Body())
		err := fn(ctx, req, rsp)
		return err
	}
}
