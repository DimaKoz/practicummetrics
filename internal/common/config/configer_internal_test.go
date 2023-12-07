package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	flag2 "github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

const (
	addressEnvName   = "ADDRESS"
	pollEnvName      = "POLL_INTERVAL"
	reportEnvName    = "REPORT_INTERVAL"
	keyEnvName       = "KEY"
	cryptoKeyEnvName = "CRYPTO_KEY"
)

var (
	longErrD = "cannot process flags variables: couldn't convert the request interval" +
		" to int, flag: abc, err: strconv.ParseInt: parsing \"abc\": invalid syntax"
	errInitConfig = errors.New(longErrD)
)

type argTestConfig struct {
	envAddress   string
	envPoll      string
	envReport    string
	envKey       string
	envCryptoKey string
	flagAddress  string
	flagKey      string
	flagPoll     string
	flagReport   string
}

var (
	//nolint:exhaustruct
	wantConfig1 = &AgentConfig{
		RateLimit: 0,
		Config: Config{
			Address: "127.0.0.1:59483", HashKey: "e",
			ConfigFile: unknownStringFieldValue,
		},
		PollInterval: 15, ReportInterval: 16,
	}
	//nolint:exhaustruct
	wantConfig4 = &AgentConfig{
		RateLimit: 0,
		Config: Config{
			Address: "127.0.0.1:59483", HashKey: "",
			ConfigFile: unknownStringFieldValue,
		},
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
			Config: Config{
				Address: "localhost:8080", HashKey: "", CryptoKey: "",
				ConfigFile: unknownStringFieldValue,
			}, RateLimit: 0,
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
			Config:    Config{Address: "127.0.0.1:59455", HashKey: "ww", ConfigFile: unknownStringFieldValue},
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
	for _, test1 := range testsCasesAgentInitConfig {
		test := test1
		t.Run(test.name, func(t *testing.T) {
			envArgsAgentInitConfig(t, addressEnvName, test.args.envAddress) // ENV setup
			envArgsAgentInitConfig(t, pollEnvName, test.args.envPoll)
			envArgsAgentInitConfig(t, reportEnvName, test.args.envReport)
			envArgsAgentInitConfig(t, keyEnvName, test.args.envKey)
			envArgsAgentInitConfig(t, cryptoKeyEnvName, test.args.envCryptoKey)
			/*			if test.args.flagAddress != "" ||
						test.args.flagKey != "" ||
						test.args.flagPoll != "" ||
						test.args.flagReport != "" { // Flags setup*/
			osArgOrig := os.Args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
			flag.CommandLine.SetOutput(io.Discard)
			os.Args = make([]string, 0)
			os.Args = append(os.Args, osArgOrig[0])
			appendArgsAgentInitConfig(&os.Args, "-a", test.args.flagAddress)
			appendArgsAgentInitConfig(&os.Args, "-k", test.args.flagKey)
			appendArgsAgentInitConfig(&os.Args, "-p", test.args.flagPoll)
			appendArgsAgentInitConfig(&os.Args, "-r", test.args.flagReport)
			appendArgsAgentInitConfig(&os.Args, "-crypto-key", test.args.envCryptoKey)
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
		Config: Config{ //nolint:exhaustruct
			Address:    defaultAddress,
			HashKey:    defaultKey,
			ConfigFile: unknownStringFieldValue,
		},
		StoreInterval:   defaultStoreInterval,
		ConnectionDB:    unknownStringFieldValue,
		TrustedSubnet:   unknownStringFieldValue,
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

func TestSetupDefaultServerValuesHasRestore(t *testing.T) {
	cfg := ServerConfig{
		Config: Config{
			Address:    "1",
			HashKey:    "2",
			CryptoKey:  "7",
			ConfigFile: "8",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		TrustedSubnet:   "6",
		hasRestore:      false,
		Restore:         false,
	}
	setupDefaultServerValues(&cfg, "", 42, "", "", true)

	assert.True(t, cfg.Restore)
}

func TestServerConfigString(t *testing.T) {
	cfg := ServerConfig{
		Config: Config{
			Address:    "1",
			HashKey:    "2",
			CryptoKey:  "7",
			ConfigFile: "8",
		},
		StoreInterval:   3,
		FileStoragePath: "4",
		ConnectionDB:    "5",
		TrustedSubnet:   "",
		hasRestore:      true,
		Restore:         true,
	}
	want := "Address: 1 \n StoreInterval: 3 \n FileStoragePath: 4 \n " +
		"ConnectionDB: 5 \n Key: 2 \n CryptoKey: 7 \n Restore: true \n"

	assert.Equal(t, want, cfg.String())
}

func TestPrepBuildValues(t *testing.T) {
	type args struct {
		bldV string
		bldD string
		bldC string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "NA",
			args: struct {
				bldV string
				bldD string
				bldC string
			}{bldV: "NA", bldD: "NA", bldC: "NA"},
			want: "Build version: NA\nBuild date: NA\nBuild commit:  NA\n",
		},
		{
			name: "Ok",
			args: struct {
				bldV string
				bldD string
				bldC string
			}{bldV: "12/22/23", bldD: "11:56:59", bldC: "fae3be770660ddc7da88309d82f25ccb74b6a25e"},
			want: "Build version: 12/22/23\nBuild date: 11:56:59\nBuild commit:  fae3be770660ddc7da88309d82f25ccb74b6a25e\n",
		},
	}
	for _, tt := range tests {
		unit := tt
		t.Run(unit.name, func(t *testing.T) {
			got := PrepBuildValues(unit.args.bldV, unit.args.bldD, unit.args.bldC)
			assert.Equalf(t, unit.want, got, "PrepBuildValues(%v, %v, %v)", unit.args.bldV, unit.args.bldD, unit.args.bldC)
		})
	}
}

func TestSetUnknownStrValue(t *testing.T) {
	target := "1"
	value := "2"
	setUnknownStrValue(&target, value)
	assert.Equal(t, value, target)
}
