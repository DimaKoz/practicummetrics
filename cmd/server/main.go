package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/labstack/echo/v4"
	flag2 "github.com/spf13/pflag"
)

func main() {
	//var address = flag2.String("abc", ":8080", ":8080 by default")
	var address string
	flag2.StringVarP(&address, "a", "a", ":8080",
		":8080 by default")
	flag2.Parse()
	e := echo.New()
	e.POST("/update/*", handler.UpdateHandler)
	e.GET("/value/*", handler.ValueHandler)
	e.GET("/", handler.RootHandler)
	fmt.Println(address)
	err := e.Start(address)
	if err != nil {
		panic(err)
	}

}
