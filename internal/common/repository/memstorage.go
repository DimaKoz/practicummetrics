package repository

import (
	"encoding/json"
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"log"
	"os"
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
// returns updated model.MetricUnit after that
func AddMetric(mu model.MetricUnit) model.MetricUnit {
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
	return mu.Clone()
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

var filePathStorage string

func SetupFilePathStorage(pFilePathStorage string) {
	filePathStorage = pFilePathStorage
}

func Load() error {
	if filePathStorage == "" {
		return fmt.Errorf("filePathStorage is empty")
	}
	data, err := os.ReadFile(filePathStorage)
	if err != nil {
		return fmt.Errorf("can't read '%s' file with error: %w", filePathStorage, err)
	}
	var m []model.MetricUnit
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("can't parse json with error: %w", err)
	}
	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	for _, v := range m {
		memStorage.storage[v.Name] = v
	}
	log.Printf("repository: loaded: %d \n", len(m))
	return nil
}

func Save() error {
	if filePathStorage == "" {
		return fmt.Errorf("filePathStorage is empty")
	}
	metrics := GetAllMetrics()
	j, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("can't marshal json with error: %w", err)
	}
	err = os.WriteFile(filePathStorage, j, 0666)
	if err != nil {
		return fmt.Errorf("can't write '%s' file with error: %w", filePathStorage, err)
	}
	log.Printf("repository: saved: %d \n", len(metrics))
	return nil
}
