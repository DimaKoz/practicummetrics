package config

import (
	"fmt"
	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

const (
	addressEnvName = "ADDRESS"
	pollEnvName    = "POLL_INTERVAL"
	reportEnvName  = "REPORT_INTERVAL"
)

func TestAgentInitConfig(t *testing.T) {
	type args struct {
		typeCfg int

		envAddress string
		envPoll    string
		envReport  string

		flagAddress string
		flagPoll    string
		flagReport  string
	}
	var tests = []struct {
		name    string
		args    args
		want    *Config
		wantErr error
	}{
		{
			name: "unknown type config",
			args: args{
				typeCfg: 3,
			},
			want:    nil,
			wantErr: fmt.Errorf("from CreateConfig: an unknown type of the config, no config for you \n"),
		},

		{
			name: "default values (agent)",
			args: args{
				typeCfg: AgentCfg,
			},
			want: &Config{
				Address:        "localhost:8080",
				PollInterval:   int64(defaultPollInterval),
				ReportInterval: int64(defaultReportInterval),
			},
			wantErr: nil,
		},
		{
			name: "env values without flags",
			args: args{
				typeCfg: AgentCfg,

				envAddress: "127.0.0.1:59483",
				envPoll:    "15",
				envReport:  "16",
			},
			want: &Config{
				Address:        "127.0.0.1:59483",
				PollInterval:   15,
				ReportInterval: 16,
			},
		},
		{
			name: "flags values without env",
			args: args{
				typeCfg:     AgentCfg,
				flagAddress: "127.0.0.1:59455",
				flagPoll:    "12",
				flagReport:  "15",
			},
			want: &Config{
				Address:        "127.0.0.1:59455",
				PollInterval:   12,
				ReportInterval: 15,
			},
		},
		{
			name: "flags values without env and without any report interval",
			args: args{
				typeCfg:     AgentCfg,
				flagAddress: "127.0.0.1:59455",
				flagPoll:    "5",
				flagReport:  "abc",
			},
			want: &Config{
				Address:        "127.0.0.1:59455",
				PollInterval:   5,
				ReportInterval: 10,
			},
		},

		{
			name: "flags values + env",
			args: args{
				typeCfg: AgentCfg,

				envAddress:  "127.0.0.1:59483",
				envPoll:     "3",
				envReport:   "4",
				flagAddress: "127.0.0.1:59455",
				flagPoll:    "12",
				flagReport:  "15",
			},
			want: &Config{
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
				err := os.Setenv(addressEnvName, tt.args.envAddress)
				newAddress := os.Getenv(addressEnvName)
				fmt.Println("new address:", newAddress, " ", err)
				t.Cleanup(func() { _ = os.Setenv(addressEnvName, origAddress) })
			}
			if tt.args.envPoll != "" {
				origPoll := os.Getenv(pollEnvName)
				_ = os.Setenv(pollEnvName, tt.args.envPoll)
				t.Cleanup(func() { _ = os.Setenv(pollEnvName, origPoll) })
			}
			if tt.args.envReport != "" {
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

			got, gotErr := CreateConfig(tt.args.typeCfg)

			assert.Equal(t, tt.wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, tt.wantErr)

			assert.Equal(t, tt.want, got, "Configs - got: %v, want: %v", got, tt.want)

		})
	}
}

func Test_processEnvError(t *testing.T) {
	wantErr := fmt.Errorf(" env parsing error: %w", fmt.Errorf("env: expected a pointer to a Struct"))
	gotErr := processEnv(nil)

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)

}

func Test_processEnvNoError(t *testing.T) {
	var wantErr error = nil
	gotErr := processEnv(&Config{})

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)

}

func Test_processEnvMock(t *testing.T) {
	flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
	flag2.CommandLine.SetOutput(io.Discard)

	osArgOrig := os.Args
	os.Args = make([]string, 0)
	os.Args = append(os.Args, osArgOrig[0])
	t.Cleanup(func() { os.Args = osArgOrig })

	old := getEnv
	defer func() { getEnv = old }()

	getEnv = func(config *Config) error {
		return fmt.Errorf("any error")
	}

	want := &Config{
		Address:        "localhost:8080",
		PollInterval:   int64(defaultPollInterval),
		ReportInterval: int64(defaultReportInterval),
	}
	var wantErr error = nil
	got, gotErr := CreateConfig(AgentCfg)

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)
	assert.Equal(t, want, got, "Configs - got: %v, want: %v", got, want)

}
