package sqldb

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	if err = createTables(conn); err != nil {
		return nil, err
	}
	sugar.Info("successfully connected to db", conn)
	sugar.Info("db:", conn)
	cfg.Restore = false

	return conn, nil
}

func createTables(pgConn *pgx.Conn) error {
	sqlString := `CREATE TABLE IF NOT EXISTS metrics
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(100) NOT NULL,
    value VARCHAR (200) NOT NULL
);`

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	if _, err := pgConn.Exec(ctx, sqlString); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}
