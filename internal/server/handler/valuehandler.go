package handler

import (
	"errors"
	"fmt"
	error2 "github.com/DimaKoz/practicummetrics/internal/common/error"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"io"
	"net/http"
	"strings"
)

const (
	okPathPartsValue = 4
	indexTypeValue   = 2
	indexNameValue   = 3
)

// ValueHandler handles `/value/`
func ValueHandler(res http.ResponseWriter, req *http.Request) {
	name, err := getNameFromPath(req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), err.StatusCode)
		return
	}
	mu := repository.GetMetricByName(name)
	if mu == nil {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err2 := io.WriteString(res, mu.Value); err2 != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Println("error for ValueHandler: ", err2)
		return
	}
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

}

func getNameFromPath(path string) (string, *error2.RequestError) {
	if path == "" {
		return "", &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}
	parts := strings.Split(path, "/")
	if len(parts) != okPathPartsValue {
		return "", &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")}
	}
	return parts[indexNameValue], nil
}
