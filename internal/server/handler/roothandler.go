package handler

import "net/http"

func RootHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}
