package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {

	cfg, err := config.CreateConfig(config.ServerCfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// from cfg:
	fmt.Println("cfg:")
	fmt.Println("address:", cfg.Address)

	e := echo.New()
	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.POST("/value/*", handler.ValueHandler)
	e.GET("/", handler.RootHandler)

	err = e.Start(cfg.Address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
