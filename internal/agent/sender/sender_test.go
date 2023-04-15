package sender

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_getUrl(t *testing.T) {
	type args struct {
		cfg *config.Config
		mu  model.MetricUnit
	}
	tests := []struct {
		name string
		args args
		want url.URL
	}{
		{
			name: "get url",
			args: args{
				cfg: &config.Config{
					Address:        "localhost:8080",
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
			want: url.URL{Scheme: "http", Host: "localhost:8080", Path: "gauge/qwerty/42.42"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getURL(tt.args.cfg, tt.args.mu)
			assert.Equal(t, tt.want, got, "getURL() = %v, want %v", got, tt.want)
		})
	}
}

func TestParcelsSend(t *testing.T) {
	type args struct {
		cfg *config.Config
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
				cfg: &config.Config{
					Address:        "localhost:8080",
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
			want: "/gauge/qwerty/42.42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				got := req.URL.Path
				_, err := rw.Write([]byte(`OK`))
				assert.NoError(t, err, "getURL() = %v, want no error", got)
				if err != nil {
					return
				}
				assert.Equal(t, tt.want, got, "getURL() = %v, want %v", got, tt.want)
				// Send response to be tested

			})
			// Start a local HTTP server
			srv := httptest.NewUnstartedServer(mock)

			// create a listener with the desired port.
			l, err := net.Listen("tcp", tt.args.cfg.Address)
			if err != nil {
				assert.NoError(t, err)
			}
			_ = srv.Listener.Close()
			srv.Listener = l

			// Start the server.
			srv.Start()

			urlUsed := getURL(tt.args.cfg, tt.args.mu)
			srv.URL = urlUsed.String()
			// Close the server when test finishes
			defer srv.Close()

			// Use Client & URL from our local test server
			ParcelsSend(tt.args.cfg, []model.MetricUnit{tt.args.mu})
		})
	}
}

/*
func readPrintedByte() {

	err := io.EOF // force an error
	if err != nil {
		ParcelsSend(&config.Config{
			Address:        "localhost:8080",
			PollInterval:   int64(2),
			ReportInterval: int64(10),
		}, []model.MetricUnit{model.MetricUnit{
			Type:       model.MetricTypeGauge,
			Name:       "qwerty",
			Value:      "42.42",
			ValueInt:   0,
			ValueFloat: 42.42,
		}})

		return
	}

}

	func captureOutput(f func()) string {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		f()
		log.SetOutput(os.Stderr)
		return buf.String()
	}
*/
func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}

func TestPrintSender(t *testing.T) {

	want := `client: could not create the request: Post "http://localhost:8080/gauge/qwerty/42.42"`
	// setting stdout to a file
	done := capture()
	ParcelsSend(&config.Config{
		Address:        "localhost:8080",
		PollInterval:   int64(2),
		ReportInterval: int64(10),
	}, []model.MetricUnit{model.MetricUnit{
		Type:       model.MetricTypeGauge,
		Name:       "qwerty",
		Value:      "42.42",
		ValueInt:   0,
		ValueFloat: 42.42,
	}})
	s, err := done()
	if err != nil {
		fmt.Printf("TestPrintSender(): err=%v, stdout=%q\n", err, s)
	}
	assert.Contains(t, s, want, "Expected %s, got %s", "this is value: test", want, s)

}
