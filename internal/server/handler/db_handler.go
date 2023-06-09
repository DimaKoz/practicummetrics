package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// BaseHandler holds *pgx.Conn.
type BaseHandler struct {
	conn *pgx.Conn
}

// NewBaseHandler returns a new BaseHandler.
func NewBaseHandler(dbConn *pgx.Conn) *BaseHandler {
	return &BaseHandler{
		conn: dbConn,
	}
}

// PingHandler handles http.StatusOK if DB alive
// or http.StatusInternalServerError if not.
func (h *BaseHandler) PingHandler(ctxEcho echo.Context) error {
	ctx := context.TODO()
	status := http.StatusInternalServerError
	var err error
	if h.conn != nil {
		if err = h.conn.Ping(ctx); err == nil {
			status = http.StatusOK
		}
	}

	if err = ctxEcho.NoContent(status); err != nil {
		err = fmt.Errorf("%w", err)
	}

	return err
}
