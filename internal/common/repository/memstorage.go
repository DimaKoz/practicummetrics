package repository

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"strconv"
	"sync"
)

var memStorageSync = &sync.Mutex{}
var memStorage = MemStorage{
	storage: make(map[string]model.MetricUnit, 0),
}

type MemStorage struct {
	storage map[string]model.MetricUnit
}

// AddMetric adds model.MetricUnit to 'memStorage.storage' storage
func AddMetric(mu model.MetricUnit) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	if mu.Type == model.MetricTypeCounter {
		found, ok := memStorage.storage[mu.Name]
		if ok {
			mu.ValueInt += found.ValueInt
			mu.Value = strconv.FormatInt(mu.ValueInt, 10)
		}
	}
	memStorage.storage[mu.Name] = mu
}

// GetMetricByName returns a model.MetricUnit and nil error if found or model.EmptyMetric and error
func GetMetricByName(name string) (model.MetricUnit, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	found, ok := memStorage.storage[name]
	if ok {
		return found.Clone(), nil
	}
	return model.EmptyMetric, fmt.Errorf("couldn't find a metric: %s", name)
}

// GetAllMetrics returns a list of model.MetricUnit from the storage
func GetAllMetrics() []model.MetricUnit {
	result := make([]model.MetricUnit, 0)
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	for _, v := range memStorage.storage {
		result = append(result, v.Clone())
	}
	return result
}
