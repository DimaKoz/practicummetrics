package handler

import "net/http"

// RootHandler handles `/`
func RootHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
}
