package config

import (
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

const (
	addressEnvName = "ADDRESS"
	pollEnvName    = "POLL_INTERVAL"
	reportEnvName  = "REPORT_INTERVAL"
)

func TestAgentInitConfig(t *testing.T) {
	type args struct {
		cfg                 *model.Config
		defaultAddress      string
		defaultRepInterval  time.Duration
		defaultPollInterval time.Duration
		envAddress          string
		envPoll             string
		envReport           string

		flagAddress string
		flagPoll    string
		flagReport  string
	}
	tests := []struct {
		name string
		args args
		want *model.Config
	}{
		{
			name: "nil cfg",
			args: args{
				cfg:                 nil,
				defaultAddress:      "localhost:8080",
				defaultPollInterval: time.Duration(2),
				defaultRepInterval:  time.Duration(10),
			},
			want: nil,
		},

		{
			name: "default values",
			args: args{
				cfg:                 &model.Config{},
				defaultAddress:      "localhost:8080",
				defaultPollInterval: time.Duration(2),
				defaultRepInterval:  time.Duration(10),
			},
			want: &model.Config{
				Address:        "localhost:8080",
				PollInterval:   2,
				ReportInterval: 10,
			},
		},
		{
			name: "env values without flags",
			args: args{
				cfg:                 &model.Config{},
				defaultAddress:      "localhost:8080",
				defaultPollInterval: time.Duration(2),
				defaultRepInterval:  time.Duration(10),
				envAddress:          "127.0.0.1:59483",
				envPoll:             "15",
				envReport:           "16",
			},
			want: &model.Config{
				Address:        "127.0.0.1:59483",
				PollInterval:   15,
				ReportInterval: 16,
			},
		},
		{
			name: "flags values without env",
			args: args{
				cfg:                 &model.Config{},
				defaultAddress:      "localhost:8080",
				defaultPollInterval: time.Duration(2),
				defaultRepInterval:  time.Duration(10),
				flagAddress:         "127.0.0.1:59455",
				flagPoll:            "12",
				flagReport:          "15",
			},
			want: &model.Config{
				Address:        "127.0.0.1:59455",
				PollInterval:   12,
				ReportInterval: 15,
			},
		},
		{
			name: "flags values + env",
			args: args{
				cfg:                 &model.Config{},
				defaultAddress:      "localhost:8080",
				defaultPollInterval: time.Duration(2),
				defaultRepInterval:  time.Duration(10),
				envAddress:          "127.0.0.1:59483",
				envPoll:             "3",
				envReport:           "4",
				flagAddress:         "127.0.0.1:59455",
				flagPoll:            "12",
				flagReport:          "15",
			},
			want: &model.Config{
				Address:        "127.0.0.1:59483",
				PollInterval:   3,
				ReportInterval: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// ENV setup
			if tt.args.envAddress != "" {
				origAddress := os.Getenv(addressEnvName)
				_ = os.Setenv(addressEnvName, tt.args.envAddress)
				t.Cleanup(func() { _ = os.Setenv(addressEnvName, origAddress) })
			}
			if tt.args.envPoll != "" {
				origPoll := os.Getenv(pollEnvName)
				_ = os.Setenv(pollEnvName, tt.args.envPoll)
				t.Cleanup(func() { _ = os.Setenv(pollEnvName, origPoll) })
			}
			if tt.args.envAddress != "" {
				origReport := os.Getenv(reportEnvName)
				_ = os.Setenv(reportEnvName, tt.args.envReport)
				t.Cleanup(func() { _ = os.Setenv(reportEnvName, origReport) })
			}

			// Flags setup

			if tt.args.flagAddress != "" || tt.args.flagPoll != "" || tt.args.flagReport != "" {
				flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
				flag2.CommandLine.SetOutput(io.Discard)

				osArgOrig := os.Args
				os.Args = make([]string, 0)
				os.Args = append(os.Args, osArgOrig[0])
				if tt.args.flagAddress != "" {
					os.Args = append(os.Args, "-a")
					os.Args = append(os.Args, tt.args.flagAddress)
				}
				if tt.args.flagPoll != "" {
					os.Args = append(os.Args, "-p")
					os.Args = append(os.Args, tt.args.flagPoll)
				}

				if tt.args.flagReport != "" {
					os.Args = append(os.Args, "-r")
					os.Args = append(os.Args, tt.args.flagReport)
				}

				t.Cleanup(func() { os.Args = osArgOrig })
			}

			got := tt.args.cfg
			AgentInitConfig(got, tt.args.defaultAddress, tt.args.defaultRepInterval, tt.args.defaultPollInterval)
			assert.Equal(t, tt.want, got, "Configs - got: %v, want: %v", got, tt.want)

		})
	}
}

func Test_processEnv(t *testing.T) {
	processEnv(nil)
}
