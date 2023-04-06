package repository

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"strconv"
	"sync"
)

var once sync.Once
var lockMemSt = &sync.Mutex{}
var instanceMemSt MemStorage

type MemStorage struct {
	storage map[string]model.MetricUnit
}

// InitMemStorage initializes the MemStorage inside itself
func InitMemStorage() {
	once.Do(func() {

		instanceMemSt = MemStorage{}
		instanceMemSt.storage = make(map[string]model.MetricUnit, 0)

	})
}

// AddMetricMemStorage adds model.MetricUnit to 'instanceMemSt.storage' storage
func AddMetricMemStorage(mu model.MetricUnit) {
	InitMemStorage()
	lockMemSt.Lock()
	defer lockMemSt.Unlock()

	if mu.Type == model.MetricTypeCounter {
		found, ok := instanceMemSt.storage[mu.Name]
		if ok {
			mu.ValueI += found.ValueI
			mu.Value = strconv.FormatInt(mu.ValueI, 10)
		}
	}
	instanceMemSt.storage[mu.Name] = mu
}

// GetMetricsMemStorage returns a list of model.MetricUnit from the storage
func GetMetricsMemStorage() []model.MetricUnit {
	result := make([]model.MetricUnit, 0)
	lockMemSt.Lock()
	defer lockMemSt.Unlock()
	for _, v := range instanceMemSt.storage {
		newMetric, err := model.NewMetricUnit(v.Type, v.Name, v.Value)
		if err == nil && newMetric != nil {
			result = append(result, *newMetric)
		}
	}
	return result
}
