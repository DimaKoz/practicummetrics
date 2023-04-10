package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
)

const defaultAddress = "localhost:8080"

func main() {
	cfg := &model.Config{}
	config.AgentInitConfig(cfg, defaultAddress, 0, 0)

	// from cfg:
	fmt.Println("cfg:")
	fmt.Println("address:", cfg.Address)
	e := echo.New()
	e.POST("/update/*", handler.UpdateHandler)
	e.GET("/value/*", handler.ValueHandler)
	e.POST("/value/*", handler.ValueHandler)
	e.GET("/", handler.RootHandler)

	err := e.Start(cfg.Address)
	if err != nil {
		panic(err)
	}

}
