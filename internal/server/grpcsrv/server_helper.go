package grpcsrv

import (
	"context"
	"fmt"
	"net"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/DimaKoz/practicummetrics/internal/server/serializer"
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
	metricsSlice := make([]model.Metrics, 0)
	if err := serializer.DeserializeString(request.Metrics.Body, &metricsSlice); err != nil {
		return &response, fmt.Errorf("can't update metrics by: %w", err)
	}
	for _, item := range metricsSlice {
		prepModelValue, err := item.GetPreparedValue()
		if err != nil {
			erDesc := fmt.Sprintf("gRPC: Metrics contains nil: %s", err)

			return &response, fmt.Errorf("%s : %w", erDesc, err)
		}
		muIncome, err := model.NewMetricUnit(item.MType, item.ID, prepModelValue)
		if err != nil {
			return &response, fmt.Errorf("gRPC: cannot create metric: %w", err)
		}
		_ = repository.AddMetric(muIncome)
	}
	s.saveUpdates()

	return &response, nil
}

// saveUpdates stores values of repository to a file.
func (s *MetricsServer) saveUpdates() {
	if s.cfg.FileStoragePath != "" && s.cfg.StoreInterval == 0 {
		go func() {
			err := repository.SaveVariant()
			if err != nil {
				zap.S().Fatal(err)
			}
		}()
	}
}
