package main

import (
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	"log"
)

func main() {

	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("couldn't create a config %s", err)
	}

	// from cfg:
	log.Println("cfg:")
	log.Println("address:", cfg.Address)

	e := echo.New()
	e.POST("/update/:type/:name/:value", handler.UpdateHandler)
	e.GET("/value/:type/:name", handler.ValueHandler)
	e.GET("/", handler.RootHandler)

	err = e.Start(cfg.Address)
	if err != nil {
		log.Fatalf("couldn't start the server by %s", err)
	}

}
