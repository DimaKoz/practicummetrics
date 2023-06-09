package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var errNoInfoConnectionDB = errors.New("no DB connection info")

// ConnectDB opens a connection to the database.
func ConnectDB(cfg *config.ServerConfig, sugar zap.SugaredLogger) (*pgx.Conn, error) {
	if cfg == nil || cfg.ConnectionDB == "" {
		return nil, errNoInfoConnectionDB
	}
	conn, err := pgx.Connect(context.Background(), cfg.ConnectionDB)
	if err != nil {
		return nil, fmt.Errorf("failed to get a DB connection: %w", err)
	}
	sugar.Info("successfully connected to db", conn)
	sugar.Info("db:", conn)

	return conn, nil
}
