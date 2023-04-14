package repository

import (
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

// GetMetricByName returns a *model.MetricUnit if found or nil
func GetMetricByName(name string) *model.MetricUnit {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	var result *model.MetricUnit = nil
	found, ok := memStorage.storage[name]
	if ok {
		result = &model.MetricUnit{
			Type:       found.Type,
			Name:       found.Name,
			Value:      found.Value,
			ValueFloat: found.ValueFloat,
			ValueInt:   found.ValueInt,
		}
	}
	return result
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
