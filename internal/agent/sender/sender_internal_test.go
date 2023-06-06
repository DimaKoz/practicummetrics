package sender

import (
	"bytes"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
					Config: config.Config{
						Address: "localhost:8181",
					},
					PollInterval:   int64(2),
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

func readByte() {
	err := io.EOF // force an error
	if err != nil {
		return
	}
}

func TestPrintSender(t *testing.T) {
	want := `Post "http://localhost:8888/update/"`

	var buf bytes.Buffer

	log.SetOutput(&buf)

	defer func() {
		log.SetOutput(os.Stderr)
	}()
	ParcelsSend(&config.AgentConfig{
		Config: config.Config{
			Address: "localhost:8888",
		},
		PollInterval:   int64(2),
		ReportInterval: int64(10),
	}, []model.MetricUnit{{
		Type:       model.MetricTypeGauge,
		Name:       "qwerty",
		Value:      "42.42",
		ValueInt:   0,
		ValueFloat: 42.42,
	}})
	readByte()
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
			assert.Equalf(t, test.want, getTargetURL(test.args.address), "getTargetURL(%v)", test.args.address)
		})
	}
}
