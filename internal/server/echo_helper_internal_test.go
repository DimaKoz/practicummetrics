package server

import (
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMiddleware(t *testing.T) {
	echoFr := echo.New()
	//nolint:exhaustruct
	cfg := config.ServerConfig{}
	SetupMiddleware(echoFr, &cfg)
	err := echoFr.Close()
	require.NoError(t, err)
}

func TestSetupRouter(t *testing.T) {
	echoFr := echo.New()
	SetupRouter(echoFr, nil)
	wantRoutePath := "/updates/"
	wantLenRoutes := 7
	router := echoFr.Router()
	routes := router.Routes()
	assert.Len(t, routes, wantLenRoutes)
	found := false
	for i := 0; i < len(routes); i++ {
		if wantRoutePath == routes[i].Path {
			found = true

			break
		}
	}
	assert.True(t, found)
	err := echoFr.Close()
	require.NoError(t, err)
}
