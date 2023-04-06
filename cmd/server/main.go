package main

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.POST("/update/*", handler.UpdateHandler)
	e.GET("/value/*", handler.ValueHandler)
	e.GET("/", handler.RootHandler)

	err := e.Start(":8080")
	if err != nil {
		panic(err)
	}

}
