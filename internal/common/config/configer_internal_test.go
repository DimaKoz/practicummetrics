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
	keyEnvName     = "KEY"
)

var (
	longErrD = "cannot process flags variables: couldn't convert the request interval" +
		" to int, flag: abc, err: strconv.ParseInt: parsing \"abc\": invalid syntax"
	errInitConfig = errors.New(longErrD)
)

type argTestConfig struct {
	envAddress  string
	envPoll     string
	envReport   string
	envKey      string
	flagAddress string
	flagKey     string
	flagPoll    string
	flagReport  string
}

var (
	wantConfig1 = &AgentConfig{
		RateLimit:    0,
		Config:       Config{Address: "127.0.0.1:59483", HashKey: "e"},
		PollInterval: 15, ReportInterval: 16,
	}
	wantConfig4 = &AgentConfig{
		RateLimit:    0,
		Config:       Config{Address: "127.0.0.1:59483", HashKey: ""},
		PollInterval: 3, ReportInterval: 4,
	}
)

//nolint:exhaustruct
var testsCasesAgentInitConfig = []struct {
	name    string
	args    argTestConfig
	want    *AgentConfig
	wantErr error
}{
	{
		name: "default values (agent)", args: argTestConfig{}, wantErr: nil, //nolint:exhaustruct
		want: &AgentConfig{
			Config: Config{Address: "localhost:8080", HashKey: ""}, RateLimit: 0,
			PollInterval: int64(defaultPollInterval), ReportInterval: int64(defaultReportInterval),
		},
	},
	{
		name: "env", args: argTestConfig{ //nolint:exhaustruct
			envAddress: "127.0.0.1:59483",
			envPoll:    "15",
			envReport:  "16",
			envKey:     "e",
		},
		want: wantConfig1,
	},
	{
		name: "flags",
		args: argTestConfig{ //nolint:exhaustruct
			flagAddress: "127.0.0.1:59455", flagPoll: "12", flagReport: "15", flagKey: "ww",
		},
		want: &AgentConfig{
			Config:    Config{Address: "127.0.0.1:59455", HashKey: "ww"},
			RateLimit: 0, PollInterval: 12, ReportInterval: 15,
		},
	},
	{
		name: "flags values without env and without any report interval", want: nil, wantErr: errInitConfig,
		args: argTestConfig{flagAddress: "127.0.0.1:59455", flagPoll: "5", flagReport: "abc"}, //nolint:exhaustruct
	},
	{
		name: "flags&env", want: wantConfig4,
		args: argTestConfig{ //nolint:exhaustruct
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
			envArgsAgentInitConfig(t, keyEnvName, test.args.envKey)
			/*			if test.args.flagAddress != "" ||
						test.args.flagKey != "" ||
						test.args.flagPoll != "" ||
						test.args.flagReport != "" { // Flags setup*/
			osArgOrig := os.Args
			flag2.CommandLine = flag2.NewFlagSet(os.Args[0], flag2.ContinueOnError)
			flag2.CommandLine.SetOutput(io.Discard)
			os.Args = make([]string, 0)
			os.Args = append(os.Args, osArgOrig[0])
			appendArgsAgentInitConfig(&os.Args, "-a", test.args.flagAddress)
			appendArgsAgentInitConfig(&os.Args, "-k", test.args.flagKey)
			appendArgsAgentInitConfig(&os.Args, "-p", test.args.flagPoll)
			appendArgsAgentInitConfig(&os.Args, "-r", test.args.flagReport)
			t.Cleanup(func() { os.Args = osArgOrig })
			//}

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
			HashKey: defaultKey,
		},
		StoreInterval:   defaultStoreInterval,
		ConnectionDB:    unknownStringFieldValue,
		FileStoragePath: defaultFileStoragePath,
		hasRestore:      true,
		Restore:         defaultRestore,
	}
	got := NewServerConfig()
	err := LoadServerConfig(got, ProcessEnvServer)
	assert.NoError(t, err, "error must be nil")
	assert.Equal(t, want, got, "Configs - got: %v, want: %v", got, want)
}

func TestServerConfigIsUseDatabase(t *testing.T) {
	tests := []struct {
		name string
		cfg  ServerConfig
		want bool
	}{
		{
			name: "use db == true",
			want: true,
			cfg: ServerConfig{ //nolint:exhaustruct
				ConnectionDB: "1234",
			},
		},
		{
			name: "use db == false, ConnectionDB is empty",
			want: false,
			cfg: ServerConfig{ //nolint:exhaustruct
				ConnectionDB: "",
			},
		},
		{
			name: "use db == false, ConnectionDB is 'unknownStringFieldValue'",
			want: false,
			cfg: ServerConfig{ //nolint:exhaustruct
				ConnectionDB: "unknownStringFieldValue",
			},
		},
	}
	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.cfg.IsUseDatabase())
		})
	}
}

func TestServerConfigString(t *testing.T) {
	cfg := ServerConfig{
		Config: Config{
			Address: "1",
			HashKey: "2",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		hasRestore:      true,
		Restore:         true,
	}
	want := "Address: 1 \n StoreInterval: 3 \n FileStoragePath: 4 \n ConnectionDB: 5 \n Key: 2 \n Restore: true \n"

	assert.Equal(t, want, cfg.String())
	assert.Equal(t, want, cfg.StringVariantBuffer())
	assert.Equal(t, want, cfg.StringVariantCopy())
}
