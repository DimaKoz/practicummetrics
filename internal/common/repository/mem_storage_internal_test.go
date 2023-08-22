package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddMetricMemStorage(t *testing.T) {
	type args struct {
		mu model.MetricUnit
	}
	tests := []struct {
		name    string
		args    []args
		wantKey string
		want    *model.MetricUnit
	}{
		{
			name: "counter",
			args: []args{
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "42", ValueInt: 42, ValueFloat: 0}},
				{mu: model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "10", ValueInt: 10, ValueFloat: 0}},
			},
			wantKey: "test",
			want:    &model.MetricUnit{Type: model.MetricTypeCounter, Name: "test", Value: "52", ValueInt: 52, ValueFloat: 0},
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			for _, unit := range test.args {
				AddMetric(unit.mu)
			}
			if got, ok := memStorage.storage[test.wantKey]; ok {
				if !reflect.DeepEqual(&got, test.want) {
					t.Errorf("AddMetric() got = %v, want %v", got, test.want)
				}
			} else {
				t.Errorf("not found stored result")
			}
		})
	}
}

var errTestGetMetricByName = errors.New("couldn't find a metric: ")

func TestGetMetricByName(t *testing.T) {
	type args struct {
		search string
		add    []model.MetricUnit
	}
	tests := []struct {
		name    string
		args    args
		want    model.MetricUnit
		wantErr error
	}{
		{
			name: "empty key",
			args: args{
				search: "",
				add:    []model.MetricUnit{},
			},
			want:    model.EmptyMetric,
			wantErr: errTestGetMetricByName,
		},
		{
			name: "wanted key",
			args: args{
				search: "wanted",
				add: []model.MetricUnit{
					{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
					{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
				},
			},
			want:    model.MetricUnit{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
			wantErr: nil,
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			orig := memStorage.storage
			memStorage.storage = make(map[string]model.MetricUnit, 0)
			t.Cleanup(func() { memStorage.storage = orig })
			for _, v := range test.args.add {
				AddMetric(v)
			}
			got, err := GetMetricByName(test.args.search)
			if test.wantErr != nil {
				assert.EqualError(t, err, test.wantErr.Error())
			} else {
				assert.Nil(t, err, " want no error")
			}

			assert.Equal(t, got, test.want, "GetMetricByName() = %v, want %v", got, test.want)
		})
	}
}

func TestGetMetricsMemStorage(t *testing.T) {
	tests := []struct {
		add  []model.MetricUnit
		name string
		want []model.MetricUnit
	}{
		{
			name: "empty",
			add:  []model.MetricUnit{},
			want: []model.MetricUnit{},
		}, {
			name: "wanted key",
			add: []model.MetricUnit{
				{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
				{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
			},
			want: []model.MetricUnit{
				{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
				{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
			},
		},
	}

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			orig := memStorage.storage
			memStorage.storage = make(map[string]model.MetricUnit, 0)
			t.Cleanup(func() { memStorage.storage = orig })
			for _, v := range test.add {
				AddMetric(v)
			}
			assert.ElementsMatch(t, test.want, GetAllMetrics(), "GetAllMetrics()")
		})
	}
}

func TestLoadSaveEmptyFileStorageErr(t *testing.T) {
	orig := filePathStorage
	filePathStorage = ""

	t.Cleanup(
		func() {
			filePathStorage = orig
		})

	err := Load()
	assert.Error(t, err)
	err = LoadVariant()
	assert.Error(t, err)

	err = Save()
	assert.Error(t, err)
	err = SaveVariant()
	assert.Error(t, err)
}

func TestSetupFilePathStorage(t *testing.T) {
	orig := filePathStorage
	want := filepath.Join(t.TempDir(), "abc.txt")

	t.Cleanup(func() { filePathStorage = orig })

	SetupFilePathStorage(want)
	assert.Equal(t, want, filePathStorage)
}

func TestLoadSave(t *testing.T) {
	zap.L()
	logger := zap.Must(zap.NewDevelopment())

	defer func(loggerZap *zap.Logger) {
		_ = loggerZap.Sync()
	}(logger)

	zap.ReplaceGlobals(logger)

	orig := filePathStorage
	filePathStorage = filepath.Join(t.TempDir(), "abc"+fmt.Sprintf("%d", time.Now().Unix())+".json")
	origMemSt := memStorage.storage
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	t.Cleanup(
		func() {
			memStorage.storage = origMemSt
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)
	want := []model.MetricUnit{
		{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
		{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
	}
	for _, v := range want {
		AddMetric(v)
	}
	err := Save()
	assert.NoError(t, err)
	memStorage.storage = make(map[string]model.MetricUnit, 0)

	mu, err := GetMetricByName("wanted")
	assert.Error(t, err)
	assert.Equal(t, model.EmptyMetric, mu)

	err = Load()
	assert.NoError(t, err)
	assert.ElementsMatch(t, want, GetAllMetrics(), "GetAllMetrics()")
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	err = LoadVariant()
	assert.NoError(t, err)
	assert.ElementsMatch(t, want, GetAllMetrics(), "GetAllMetrics()")
}

func TestLoadSaveVariant(t *testing.T) {
	orig := filePathStorage
	filePathStorage = filepath.Join(t.TempDir(), "test"+fmt.Sprintf("%d", time.Now().Unix())+".json")
	origMemSt := memStorage.storage
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	t.Cleanup(
		func() {
			memStorage.storage = origMemSt
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)
	want := []model.MetricUnit{
		{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
		{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
	}
	for _, v := range want {
		AddMetric(v)
	}
	err := SaveVariant()
	assert.NoError(t, err)
	memStorage.storage = make(map[string]model.MetricUnit, 0)

	mu, err := GetMetricByName("wanted")
	assert.Error(t, err)
	assert.Equal(t, model.EmptyMetric, mu)

	err = LoadVariant()
	assert.NoError(t, err)
	assert.ElementsMatch(t, want, GetAllMetrics(), "GetAllMetrics()")
}

func TestLoadErrorFile(t *testing.T) {
	orig := filePathStorage
	filePathStorage = filepath.Join(t.TempDir(), "abc"+fmt.Sprintf("%d", time.Now().Unix())+".json")

	t.Cleanup(
		func() {
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)

	err := Load()
	assert.Error(t, err)
	err = LoadVariant()
	assert.Error(t, err)
}

func TestLoadErrorParse(t *testing.T) {
	orig := filePathStorage
	filePathStorage = filepath.Join(t.TempDir(), "test"+fmt.Sprintf("%d", time.Now().Unix())+".json")
	var perm os.FileMode = 0o600
	err := os.WriteFile(filePathStorage, []byte{'{'}, perm)
	require.NoError(t, err)
	t.Cleanup(
		func() {
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)

	err = Load()
	assert.Error(t, err)
}

/*

BenchmarkLoad
BenchmarkLoad/Load()
BenchmarkLoad/Load()-8         	   54238	     21059 ns/op	    1512 B/op	      21 allocs/op
BenchmarkLoad/LoadVariant()
BenchmarkLoad/LoadVariant()-8  	   83830	     14361 ns/op	    1384 B/op	      12 allocs/op
BenchmarkSave
BenchmarkSave/Save()
BenchmarkSave/Save()-8         	    4788	    220934 ns/op	     248 B/op	       5 allocs/op
BenchmarkSave/SaveVariant()
BenchmarkSave/SaveVariant()-8  	    4921	    239360 ns/op	     243 B/op	       5 allocs/op

*/

func BenchmarkLoad(b *testing.B) {
	orig := filePathStorage
	filePathStorage = filepath.Join(b.TempDir(), "abc"+fmt.Sprintf("%d", time.Now().Unix())+".json")
	origMemSt := memStorage.storage
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.Cleanup(
		func() {
			memStorage.storage = origMemSt
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)
	want := []model.MetricUnit{
		{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
		{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
	}
	for _, v := range want {
		AddMetric(v)
	}
	err := Save()
	require.NoError(b, err)
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.ResetTimer()
	b.Run("Load()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Load()
		}
	})
	b.StopTimer()
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.StartTimer()
	b.Run("LoadVariant()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = LoadVariant()
		}
	})
}

func BenchmarkSave(b *testing.B) {
	orig := filePathStorage
	filePathStorage = filepath.Join(b.TempDir(), "bench"+fmt.Sprintf("%d", time.Now().Unix())+".json")
	origMemSt := memStorage.storage
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.Cleanup(
		func() {
			memStorage.storage = origMemSt
			filePathStorage = orig
		})
	SetupFilePathStorage(filePathStorage)
	want := []model.MetricUnit{
		{Type: model.MetricTypeCounter, Name: "wanted", Value: "42", ValueInt: 42, ValueFloat: 0},
		{Type: model.MetricTypeCounter, Name: "not_wanted", Value: "43", ValueInt: 43, ValueFloat: 0},
		{Type: model.MetricTypeCounter, Name: "not_wanted1", Value: "44", ValueInt: 44, ValueFloat: 0},
	}
	for _, v := range want {
		AddMetric(v)
	}

	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.ResetTimer()
	b.Run("Save()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Save()
		}
	})
	b.StopTimer()
	memStorage.storage = make(map[string]model.MetricUnit, 0)
	b.StartTimer()
	b.Run("SaveVariant()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = SaveVariant()
		}
	})
}
