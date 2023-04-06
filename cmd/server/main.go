package main

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	//gin  emulated
	r := gin.Default()
	fmt.Println(r)

	mux := http.NewServeMux()
	mux.Handle(`/update/`, middleware.MethodChecker(http.HandlerFunc(handler.UpdateHandler)))
	mux.Handle(`/value/`, http.HandlerFunc(handler.ValueHandler))
	mux.HandleFunc(`/`, handler.RootHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
