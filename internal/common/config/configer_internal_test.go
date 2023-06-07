package config

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

const (
	addressEnvName = "ADDRESS"
	pollEnvName    = "POLL_INTERVAL"
	reportEnvName  = "REPORT_INTERVAL"
)

var (
	longErrD = "cannot process flags variables: couldn't convert the request interval" +
		" to int, rFlag: abc, err: strconv.ParseInt: parsing \"abc\": invalid syntax"
	errInitConfig = errors.New(longErrD)
)

type argTestConfig struct {
	envAddress  string
	envPoll     string
	envReport   string
	flagAddress string
	flagPoll    string
	flagReport  string
}

var (
	wantConfig1 = &AgentConfig{Config: Config{Address: "127.0.0.1:59483"}, PollInterval: 15, ReportInterval: 16}
	wantConfig4 = &AgentConfig{Config: Config{Address: "127.0.0.1:59483"}, PollInterval: 3, ReportInterval: 4}
)

var testsCasesAgentInitConfig = []struct {
	name    string
	args    argTestConfig
	want    *AgentConfig
	wantErr error
}{
	{
		name: "default values (agent)", args: argTestConfig{}, wantErr: nil, //nolint:exhaustruct
		want: &AgentConfig{
			Config:       Config{Address: "localhost:8080"},
			PollInterval: int64(defaultPollInterval), ReportInterval: int64(defaultReportInterval),
		},
	},
	{
		name: "env", args: argTestConfig{envAddress: "127.0.0.1:59483", envPoll: "15", envReport: "16"}, //nolint:exhaustruct
		want: wantConfig1,
	},
	{
		name: "flags",
		args: argTestConfig{flagAddress: "127.0.0.1:59455", flagPoll: "12", flagReport: "15"}, //nolint:exhaustruct
		want: &AgentConfig{Config: Config{Address: "127.0.0.1:59455"}, PollInterval: 12, ReportInterval: 15},
	},
	{
		name: "flags values without env and without any report interval", want: nil, wantErr: errInitConfig,
		args: argTestConfig{flagAddress: "127.0.0.1:59455", flagPoll: "5", flagReport: "abc"}, //nolint:exhaustruct
	},
	{
		name: "flags&env", want: wantConfig4,
		args: argTestConfig{
			envAddress: "127.0.0.1:59483", envPoll: "3", envReport: "4",
			flagAddress: "127.0.0.1:59455", flagPoll: "12", flagReport: "15",
		},
	},
}

func TestAgentInitConfig(t *testing.T) {
	for _, test := range testsCasesAgentInitConfig {
		test := test
		t.Run(test.name, func(t *testing.T) {
			envArgsAgentInitConfig(t, addressEnvName, test.args.envAddress) // ENV setup
			envArgsAgentInitConfig(t, pollEnvName, test.args.envPoll)
			envArgsAgentInitConfig(t, reportEnvName, test.args.envReport)
			if test.args.flagAddress != "" || test.args.flagPoll != "" || test.args.flagReport != "" { // Flags setup
				osArgOrig := os.Args
				flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
				flag2.CommandLine.SetOutput(io.Discard)
				os.Args = make([]string, 0)
				os.Args = append(os.Args, osArgOrig[0])
				appendArgsAgentInitConfig(&os.Args, "-a", test.args.flagAddress)
				appendArgsAgentInitConfig(&os.Args, "-p", test.args.flagPoll)
				appendArgsAgentInitConfig(&os.Args, "-r", test.args.flagReport)
				t.Cleanup(func() { os.Args = osArgOrig })
			}

			got, gotErr := LoadAgentConfig()

			if test.wantErr != nil {
				assert.EqualErrorf(t, gotErr, test.wantErr.Error(), "Configs - got error: %v, want: %v", gotErr, test.wantErr)
			} else {
				assert.NoError(t, gotErr, "Configs - got error: %v, want: %v", gotErr, test.wantErr)
			}

			assert.Equal(t, test.want, got, "Configs - got: %v, want: %v", got, test.want)
		})
	}
}

func envArgsAgentInitConfig(t *testing.T, key string, value string) {
	t.Helper()
	if value != "" {
		origValue := os.Getenv(key)
		err := os.Setenv(key, value)
		log.Println("new "+key+":", value, " err:", err)
		assert.NoError(t, err)
		t.Cleanup(func() { _ = os.Setenv(key, origValue) })
	}
}

func appendArgsAgentInitConfig(target *[]string, key string, value string) {
	if value != "" {
		*target = append(*target, key)
		*target = append(*target, value)
	}
}

var errTestProcessEnvError = errors.New("env: expected a pointer to a Struct")

func TestProcessEnvError(t *testing.T) {
	wantErr := fmt.Errorf("failed to parse an environment, error: %w", errTestProcessEnvError)
	gotErr := ProcessEnvServer(nil)

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)
}

func TestProcessEnvNoError(t *testing.T) {
	var wantErr error
	gotErr := ProcessEnvServer(NewServerConfig())

	assert.Equal(t, wantErr, gotErr, "Configs - got error: %v, want: %v", gotErr, wantErr)
}

var (
	errAny                    = errors.New("any error")
	errWantTestProcessEnvMock = errors.New("server config: cannot process ENV variables: any error")
)

func TestProcessEnvMock(t *testing.T) {
	flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
	flag2.CommandLine.SetOutput(io.Discard)

	osArgOrig := os.Args
	os.Args = make([]string, 0)
	os.Args = append(os.Args, osArgOrig[0])

	t.Cleanup(func() {
		os.Args = osArgOrig
	})

	processEnv := func(config *ServerConfig) error {
		return errAny
	}

	want := NewServerConfig()
	got := NewServerConfig()
	gotErr := LoadServerConfig(got, processEnv)
	wantErr := errWantTestProcessEnvMock
	assert.Equal(t, wantErr.Error(), gotErr.Error(), "Configs - got error: %v, want: %v", gotErr, wantErr)
	assert.Equal(t, want, got, "Configs - got: %v, want: %v", got, want)
}

func TestLoadServerConfig(t *testing.T) {
	want := &ServerConfig{
		Config: Config{
			Address: defaultAddress,
		},
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		hasRestore:      true,
		Restore:         defaultRestore,
	}
	got := NewServerConfig()
	err := LoadServerConfig(got, ProcessEnvServer)
	assert.NoError(t, err, "error must be nil")
	assert.Equal(t, want, got, "Configs - got: %v, want: %v", got, want)
}
