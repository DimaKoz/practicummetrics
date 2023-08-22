package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	goccyj "github.com/goccy/go-json"
	"go.uber.org/zap"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
)

var (
	memStorageSync = &sync.Mutex{}
	memStorage     = MemStorage{
		storage: make(map[string]model.MetricUnit, 0),
	}
)

type MemStorage struct {
	storage map[string]model.MetricUnit
}

var errRepo = errors.New("couldn't find a metric")

func repositoryError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

// AddMetric adds model.MetricUnit to '_memStorage.storage' storage
// returns updated model.MetricUnit after that.
func AddMetric(metricUnit model.MetricUnit) model.MetricUnit {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	if metricUnit.Type == model.MetricTypeCounter {
		found, ok := memStorage.storage[metricUnit.Name]
		if ok {
			metricUnit.ValueInt += found.ValueInt
			metricUnit.Value = strconv.FormatInt(metricUnit.ValueInt, 10)
		}
	}
	memStorage.storage[metricUnit.Name] = metricUnit

	return metricUnit.Clone()
}

// GetMetricByName returns a model.MetricUnit and nil error if found or model.EmptyMetric and error.
func GetMetricByName(name string) (model.MetricUnit, error) {
	memStorageSync.Lock()
	defer memStorageSync.Unlock()

	if found, ok := memStorage.storage[name]; ok {
		return found.Clone(), nil
	}

	return model.EmptyMetric, repositoryError(errRepo, name)
}

// GetAllMetrics returns a list of model.MetricUnit from the storage.
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

var errEmptyPath = errors.New("filePathStorage is empty")

func Load() error {
	var metricUnits []model.MetricUnit

	if filePathStorage == "" {
		return errEmptyPath
	}
	data, err := os.ReadFile(filePathStorage)
	if err != nil {
		return fmt.Errorf("can't read '%s' file with error: %w", filePathStorage, err)
	}

	if err = json.Unmarshal(data, &metricUnits); err != nil {
		return fmt.Errorf("failed to parse json with error: %w", err)
	}

	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	for _, v := range metricUnits {
		memStorage.storage[v.Name] = v
	}

	zap.S().Infof("repository: loaded: %d \n", len(metricUnits))

	return nil
}

func LoadVariant() error {
	var metricUnits []model.MetricUnit

	if filePathStorage == "" {
		return errEmptyPath
	}
	file, err := os.Open(filePathStorage)
	if err != nil {
		return fmt.Errorf("can't read '%s' file with error: %w", filePathStorage, err)
	}
	defer file.Close()
	const bufferSize = 128
	r := bufio.NewReaderSize(file, bufferSize)
	dec := goccyj.NewDecoder(r)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		return fmt.Errorf("failed to parse json with error: %w", err)
	}
	// while the array contains values

	for dec.More() {
		var mUnit model.MetricUnit

		err = dec.Decode(&mUnit)
		if err != nil {
			return fmt.Errorf("failed to parse json with error: %w", err)
		}

		metricUnits = append(metricUnits, mUnit)
	}

	// read closing bracket
	//  t, err = dec.Token()

	memStorageSync.Lock()
	defer memStorageSync.Unlock()
	for _, v := range metricUnits {
		memStorage.storage[v.Name] = v
	}

	zap.S().Infof("repository: loaded: %d \n", len(metricUnits))

	return nil
}

func Save() error {
	if filePathStorage == "" {
		return errEmptyPath
	}

	metrics := GetAllMetrics()

	var (
		saviningJSON []byte
		err          error
	)
	if saviningJSON, err = json.Marshal(metrics); err != nil {
		return fmt.Errorf("can't marshal json with error: %w", err)
	}
	var perm os.FileMode = 0o600
	if err = os.WriteFile(filePathStorage, saviningJSON, perm); err != nil {
		return fmt.Errorf("can't write '%s' file with error: %w", filePathStorage, err)
	}

	zap.S().Infof("repository: saved: %d \n", len(metrics))

	return nil
}

func SaveVariant() error {
	if filePathStorage == "" {
		return errEmptyPath
	}

	metrics := GetAllMetrics()

	var (
		saviningJSON []byte
		err          error
	)
	if saviningJSON, err = goccyj.Marshal(metrics); err != nil {
		return fmt.Errorf("can't marshal json with error: %w", err)
	}
	var perm os.FileMode = 0o600
	if err = os.WriteFile(filePathStorage, saviningJSON, perm); err != nil {
		return fmt.Errorf("can't write '%s' file with error: %w", filePathStorage, err)
	}

	zap.S().Infof("repository: saved: %d \n", len(metrics))

	return nil
}
