package grpc

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcSender struct {
	conn   *grpc.ClientConn
	client proto.MetricsClient
	logger *log.Logger
	cfg    config.AgentConfig
}

var (
	grpcSenderSync     = &sync.Mutex{}
	grpcSenderInstance grpcSender
)

func Init(cfg config.AgentConfig, logger *log.Logger) error {
	grpcSenderSync.Lock()
	defer grpcSenderSync.Unlock()
	connGrpc, err := grpc.Dial("localhost:3201", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("can't start gRPC client by: %w", err)
	}
	clientGrpc := proto.NewMetricsClient(connGrpc)
	grpcSenderInstance = grpcSender{
		conn:   connGrpc,
		logger: logger,
		cfg:    cfg,
		client: clientGrpc,
	}

	return nil
}

func Send(ctx context.Context, body string) {
	grpcSenderSync.Lock()
	defer grpcSenderSync.Unlock()
	if grpcSenderInstance.conn == nil {
		return
	}

	mp := &proto.MetricsProto{
		Body: body,
	}
	updR := &proto.UpdateRequest{
		Metrics: mp,
	}
	_, err := grpcSenderInstance.client.Updates(ctx, updR)
	if err != nil {
		grpcSenderInstance.logger.Println(err)
	}
}

func Close() {
	grpcSenderSync.Lock()
	defer grpcSenderSync.Unlock()
	if grpcSenderInstance.conn != nil {
		_ = grpcSenderInstance.conn.Close()
		grpcSenderInstance.conn = nil
	}
}
