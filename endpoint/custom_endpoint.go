package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/muhammadisa/mimir/middleware"
)

type EpCriteria struct {
	EP         endpoint.Endpoint
	DecodeFunc func(_ context.Context, request interface{}) (interface{}, error)
	EncodeFunc func(_ context.Context, response interface{}) (interface{}, error)
}

func (ec EpCriteria) AttachDefaultMiddleware(cbCommand, endpointName string, logger log.Logger) {
	const keyVals = `endpoint`
	ec.EP = middleware.LoggingMiddleware(log.With(logger, keyVals, endpointName))(ec.EP)
	ec.EP = middleware.CircuitBreakerMiddleware(cbCommand)(ec.EP)
}

type GRPCServers map[string]grpc.Handler

type Endpoints struct {
	Ep map[string]EpCriteria
}

func (e Endpoints) ApplyMiddlewares(serviceName string, logger log.Logger) {
	for endpointName := range e.Ep {
		e.Ep[endpointName].AttachDefaultMiddleware(serviceName, endpointName, logger)
	}
}

func (e Endpoints) GenerateGRPCServers(serverOptions []grpc.ServerOption) GRPCServers {
	grpcServers := make(map[string]grpc.Handler)
	for k, v := range e.Ep {
		grpcServers[k] = grpc.NewServer(v.EP, v.DecodeFunc, v.EncodeFunc, serverOptions...)
	}
	return grpcServers
}
