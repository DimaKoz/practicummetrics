package grpcsrv

import (
	"context"
	"fmt"
	proto2 "github.com/DimaKoz/practicummetrics/pkg/proto"
	"net"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MetricsServer поддерживает все необходимые методы сервера.
type MetricsServer struct {
	grpc.Server
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	proto2.UnimplementedMetricsServer

	// используем config.ServerConfig для хранения настроек
	cfg config.ServerConfig
}

func New(cfg config.ServerConfig) (*MetricsServer, error) {
	s := &MetricsServer{
		Server: *grpc.NewServer(),
		cfg:    cfg,
	}

	// регистрируем сервис
	proto2.RegisterMetricsServer(s, s)

	return s, nil
}

func (s *MetricsServer) Run(_ context.Context) error {
	// ToDo: address from cfg
	listen, err := net.Listen("tcp", ":3201")
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
