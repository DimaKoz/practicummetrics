package sender

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/go-resty/resty/v2"
	goccyj "github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestParcelsSend(t *testing.T) {
	type args struct {
		cfg *config.AgentConfig
		mu  model.MetricUnit
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check sending",
			args: args{
				cfg: &config.AgentConfig{
					Config: config.Config{ //nolint:exhaustruct
						Address: "localhost:8181", HashKey: "asdf",
					},
					RateLimit: 0, PollInterval: int64(2),
					ReportInterval: int64(10),
				},
				mu: model.MetricUnit{
					Type:       model.MetricTypeGauge,
					Name:       "qwerty",
					Value:      "42.42",
					ValueInt:   0,
					ValueFloat: 42.42,
				},
			},
			want: "{\"id\":\"qwerty\",\"type\":\"gauge\",\"value\":42.42}",
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			mock := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				// Test request parameters
				defer req.Body.Close()
				body, err := io.ReadAll(req.Body)
				assert.NoError(t, err, "want no error")
				got := string(body)
				_, err = responseWriter.Write([]byte(`OK`))
				assert.NoError(t, err, "got: %v, want no error", got)
				assert.Equal(t, test.want, got, "got: %v, want: %v", got, test.want)
			})
			// Start a local HTTP server
			srv := httptest.NewUnstartedServer(mock)

			// create a listener with the desired port.
			l, err := net.Listen("tcp", test.args.cfg.Address)
			assert.NoError(t, err)

			_ = srv.Listener.Close()
			srv.Listener = l
			srv.Start() // Start the server.

			defer srv.Close() // Close the server when test finishes

			// Use Client & URL from our local test server
			ParcelsSend(test.args.cfg, []model.MetricUnit{test.args.mu})
		})
	}
}

func TestParcelsSendBatch(t *testing.T) {
	type args struct {
		cfg *config.AgentConfig
		mu0 model.MetricUnit
		mu1 model.MetricUnit
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check sending",
			args: args{
				cfg: &config.AgentConfig{
					Config:    config.Config{Address: "localhost:8181", HashKey: "werr"}, //nolint:exhaustruct
					RateLimit: 0, PollInterval: int64(2),
					ReportInterval: int64(10),
				},
				mu0: model.MetricUnit{Type: model.MetricTypeGauge, Name: "qwerty", Value: "42.42", ValueInt: 0, ValueFloat: 42.42},
				mu1: model.MetricUnit{Type: model.MetricTypeCounter, Name: "counter", Value: "1", ValueInt: 1, ValueFloat: 0},
			},
			want: "[{\"id\":\"qwerty\",\"type\":\"gauge\",\"value\":42.42}," +
				"{\"id\":\"counter\",\"type\":\"counter\",\"delta\":1}]",
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			mock := http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
				// Test request parameters
				defer req.Body.Close()
				body, err := io.ReadAll(req.Body)
				assert.NoError(t, err, "want no error")
				got := string(body)
				responseWriter.WriteHeader(http.StatusBadRequest)
				assert.NoError(t, err, "got: %v, want no error", got)
				assert.Equal(t, test.want, got, "got: %v, want: %v", got, test.want)
			})
			// Start a local HTTP server
			srv := httptest.NewUnstartedServer(mock)

			// create a listener with the desired port.
			l, err := net.Listen("tcp", test.args.cfg.Address)
			assert.NoError(t, err)

			_ = srv.Listener.Close()
			srv.Listener = l
			srv.Start() // Start the server.

			defer srv.Close() // Close the server when test finishes

			// Use Client & URL from our local test server
			ParcelsSend(test.args.cfg, []model.MetricUnit{test.args.mu0, test.args.mu1})
		})
	}
}

func newCustomLogger(pipeTo io.Writer) zapcore.Core { //nolint:ireturn
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zap.CombineWriteSyncers(os.Stderr, zapcore.AddSync(pipeTo)),
		zapcore.InfoLevel,
	)
}

func TestPrintSender(t *testing.T) {
	want := "\\\"http://localhost:8888/update/"

	buf := &bytes.Buffer{}
	// call the constructor from your test code with the arbitrary writer
	mycore := newCustomLogger(buf)

	logger := zap.New(mycore)
	zap.ReplaceGlobals(logger)
	ParcelsSend(&config.AgentConfig{
		Config: config.Config{ //nolint:exhaustruct
			Address: "localhost:8888", HashKey: "",
		},
		RateLimit:      0,
		PollInterval:   int64(2),
		ReportInterval: int64(10),
	}, []model.MetricUnit{{
		Type:       model.MetricTypeGauge,
		Name:       "qwerty",
		Value:      "42.42",
		ValueInt:   0,
		ValueFloat: 42.42,
	}})
	got := buf.String()
	assert.Contains(t, got, want, "Expected %s, got %s", want, got)
}

func TestBrokenBatchCoverage(t *testing.T) {
	want := "http://localhost:8888/updates/"

	buf := &bytes.Buffer{}
	// call the constructor from your test code with the arbitrary writer
	mycore := newCustomLogger(buf)

	logger := zap.New(mycore)
	zap.ReplaceGlobals(logger)
	ParcelsSend(&config.AgentConfig{
		Config: config.Config{ //nolint:exhaustruct
			Address: "localhost:8888", HashKey: "",
		},
		RateLimit:      0,
		PollInterval:   int64(2),
		ReportInterval: int64(10),
	}, []model.MetricUnit{{
		Type:       model.MetricTypeGauge,
		Name:       "qwerty",
		Value:      "42.42",
		ValueInt:   0,
		ValueFloat: 42.42,
	}, {
		Type:       model.MetricTypeGauge,
		Name:       "qwerty1",
		Value:      "0.1",
		ValueInt:   0,
		ValueFloat: 0.1,
	}})
	got := buf.String()
	assert.Contains(t, got, want, "Expected %s, got %s", want, got)
}

func TestGetTargetURL(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test ok",
			args: args{
				address: "example.com",
			},
			want: "http://example.com/update/",
		},
	}
	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			got := getMetricsUpdateTargetURL(test.args.address, endpointParcelSend)
			assert.Equalf(t, test.want, got, "getMetricsUpdateTargetURL(%v)", test.args.address)
		})
	}
}

func TestAppendHashOtherMarshaling(t *testing.T) {
	metUnit, _ := model.NewMetricUnit(model.MetricTypeGauge, "RandomValue", "4321")
	emptyMetrics := model.NewEmptyMetrics()
	emptyMetrics.UpdateByMetricUnit(metUnit)
	request := resty.New().R()
	body, err := goccyj.Marshal(emptyMetrics)
	require.NoError(t, err)

	want := "dbbf6accbe592dd8b30faf871f689bc70d0e7760c0063232394f78683eac98cd" //nolint:nolintlint,gosec

	appendHashOtherMarshaling(request, "hash", body)
	got := request.Header.Get(common.HashKeyHeaderName)
	assert.Equal(t, want, got)
}
