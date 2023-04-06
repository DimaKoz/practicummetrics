package main

import (
	"github.com/DimaKoz/practicummetrics/internal/server/handler"
	"github.com/DimaKoz/practicummetrics/internal/server/middleware"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle(`/update/`, middleware.MethodChecker(http.HandlerFunc(handler.UpdateHandler)))
	mux.Handle(`/value/`, http.HandlerFunc(handler.ValueHandler))
	mux.HandleFunc(`/`, handler.RootHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}

}
