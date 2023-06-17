package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Prepare(context.Context, string, string) (*pgconn.StatementDescription, error)
	Close(context.Context) error
}

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
	timeout := 10
	var db PgxIface = conn
	if err = createTables(&db, timeout); err != nil {
		return nil, err
	}
	sugar.Info("successfully connected to db", conn)
	sugar.Info("db:", conn)
	cfg.Restore = false

	return conn, nil
}

func createTables(pgConn *PgxIface, timeout int) error {
	sqlString := `
CREATE TABLE IF NOT EXISTS metrics
(
    id    SERIAL PRIMARY KEY,
    name  VARCHAR(200) NOT NULL,
    type  VARCHAR(100) NOT NULL,
    value VARCHAR(200) NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_metrics_name
    ON metrics USING hash (name);
`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	if _, err := (*pgConn).Exec(ctx, sqlString); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
