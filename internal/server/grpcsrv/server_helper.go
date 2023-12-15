package grpcsrv

import (
	"context"
	"fmt"
	"net"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	proto2 "github.com/DimaKoz/practicummetrics/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MetricsServer supports all required methods of the 'gRPC' server.
type MetricsServer struct {
	grpc.Server
	// UnimplementedMetricsServer must be embedded to have future compatible implementations.
	proto2.UnimplementedMetricsServer

	cfg config.ServerConfig
}

func New(cfg config.ServerConfig) (*MetricsServer, error) {
	serverGrpc := &MetricsServer{
		Server:                     *grpc.NewServer(),
		UnimplementedMetricsServer: proto2.UnimplementedMetricsServer{},
		cfg:                        cfg,
	}

	proto2.RegisterMetricsServer(serverGrpc, serverGrpc)

	return serverGrpc, nil
}

func (s *MetricsServer) Run(_ context.Context) error {
	listen, err := net.Listen("tcp", "localhost:3201")
	if err != nil {
		zap.S().Warn(err)

		return fmt.Errorf("can't start gRPC server by: %w", err)
	}
	zap.S().Info("Сервер gRPC начал работу")
	go func() {
		if err = s.Serve(listen); err != nil {
			zap.S().Warn(err)

			return
		}
	}()

	return nil
}

func (s *MetricsServer) Updates(_ context.Context, request *proto2.UpdateRequest) (*emptypb.Empty, error) {
	var response emptypb.Empty
	zap.S().Info("gRPC UpdateRequest:")
	zap.S().Info(request.Metrics.Body)

	return &response, nil
}
