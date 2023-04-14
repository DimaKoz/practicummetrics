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

// AddMetricMemStorage adds model.MetricUnit to 'memStorage.storage' storage
func AddMetricMemStorage(mu model.MetricUnit) {
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

// GetMetricsMemStorage returns a list of model.MetricUnit from the storage
func GetMetricsMemStorage() []model.MetricUnit {
	result := make([]model.MetricUnit, 0)
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	for _, v := range memStorage.storage {
		newMetric, err := model.NewMetricUnit(v.Type, v.Name, v.Value)
		if err == nil && newMetric != nil {
			result = append(result, *newMetric)
		}
	}
	return result
}
