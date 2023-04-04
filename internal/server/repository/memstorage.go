package repository

import (
	"github.com/DimaKoz/practicummetrics/internal/server/model"
	"strconv"
	"sync"
)

var once sync.Once
var lockMemSt = &sync.Mutex{}
var instanceMemSt MemStorage

type MemStorage struct {
	storage map[string]model.MetricUnit
}

func InitMemStorage() {
	once.Do(func() {

		instanceMemSt = MemStorage{}
		instanceMemSt.storage = make(map[string]model.MetricUnit, 0)

	})
}

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
