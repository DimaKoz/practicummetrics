package handler

import (
	"errors"
	error2 "github.com/DimaKoz/practicummetrics/internal/server/error"
	"github.com/DimaKoz/practicummetrics/internal/server/model"
	"github.com/DimaKoz/practicummetrics/internal/server/repository"
	"net/http"
	"strings"
)

const (
	okPathParts = 5
	indexType   = 2
	indexName   = 3
	indexValue  = 4
)

// UpdateHandler handles `/update/`
func UpdateHandler(res http.ResponseWriter, req *http.Request) {
	mu, err := processPath(req.URL.Path)
	if err != nil {
		http.Error(res, err.Error(), err.StatusCode)
		return
	}
	repository.AddMetricMemStorage(*mu)
	res.WriteHeader(http.StatusOK)
}

func processPath(path string) (*model.MetricUnit, *error2.RequestError) {

	if path == "" {
		return nil, &error2.RequestError{StatusCode: http.StatusBadRequest, Err: errors.New("unavailable")}
	}

	parts := strings.Split(path, "/")
	if len(parts) != okPathParts {
		return nil, &error2.RequestError{StatusCode: http.StatusNotFound, Err: errors.New("wrong number of the parts of the path")}
	}

	mu, err := model.NewMetricUnit(parts[indexType], parts[indexName], parts[indexValue])
	if err != nil {
		return nil, err
	}

	return mu, nil

}
