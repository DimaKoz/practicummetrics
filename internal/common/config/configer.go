package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	flag2 "github.com/spf13/pflag"
)

const (
	defaultPollInterval   = time.Duration(2)
	defaultReportInterval = time.Duration(10)
	defaultAddress        = "localhost:8080"

	defaultStoreInterval    = 300
	unknownIntFieldValue    = -1
	defaultFileStoragePath  = "/tmp/metrics-db.json"
	unknownStringFieldValue = "unknownStringFieldValue"
	defaultKey              = ""
	defaultRestore          = true
)

// Config represents a config of the agent and/or the server.
type Config struct {
	Address string `env:"ADDRESS"`
	HashKey string `env:"KEY"`
}

// AgentConfig represents a config of the agent.
type AgentConfig struct {
	Config
	ReportInterval int64 `env:"REPORT_INTERVAL"`
	PollInterval   int64 `env:"POLL_INTERVAL"`
	RateLimit      int64 `env:"RATE_LIMIT"`
}

// ServerConfig represents a config of the server.
type ServerConfig struct {
	Config
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	ConnectionDB    string `env:"DATABASE_DSN"`
	hasRestore      bool
	Restore         bool `env:"RESTORE"`
}

// NewServerConfig creates an instance of ServerConfig.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Config:          Config{Address: unknownStringFieldValue, HashKey: unknownStringFieldValue},
		StoreInterval:   unknownIntFieldValue,
		FileStoragePath: unknownStringFieldValue,
		ConnectionDB:    unknownStringFieldValue,
		hasRestore:      false,
		Restore:         true,
	}
}

// ProcessEnv receives and sets up the ServerConfig.
type ProcessEnv func(config *ServerConfig) error

// LoadServerConfig loads data to the passed ServerConfig.
func LoadServerConfig(cfg *ServerConfig, processing ProcessEnv) error {
	if err := processing(cfg); err != nil {
		return fmt.Errorf("server config: cannot process ENV variables: %w", err)
	}

	if err := processServerFlags(cfg); err != nil {
		return fmt.Errorf("server config: cannot process flags variables: %w", err)
	}

	setupDefaultServerValues(cfg,
		defaultAddress,
		defaultStoreInterval,
		defaultFileStoragePath,
		defaultKey,
		defaultRestore)

	return nil
}

// LoadAgentConfig returns *AgentConfig.
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{} //nolint:exhaustruct
	cfg.HashKey = unknownStringFieldValue

	if err := processEnvAgent(cfg); err != nil {
		return nil, fmt.Errorf("agent config: cannot process ENV variables: %w", err)
	}

	if err := processAgentFlags(cfg); err != nil {
		return nil, fmt.Errorf("cannot process flags variables: %w", err)
	}

	setupDefaultAgentValues(cfg, defaultAddress, defaultReportInterval, defaultPollInterval)

	return cfg, nil
}

func addServerFlags(cfg *ServerConfig,
	address *string, rFlag *string, iFlag *string, fFlag *string, dFlag *string, keyFlag *string,
) {
	if cfg.Address == unknownStringFieldValue {
		flag2.StringVarP(address, "a", "a", unknownStringFieldValue, "")
	}

	if cfg.HashKey == unknownStringFieldValue {
		flag2.StringVarP(keyFlag, "k", "k", unknownStringFieldValue, "")
	}

	if !cfg.hasRestore {
		flag2.StringVarP(rFlag, "r", "r", unknownStringFieldValue, "")
	}

	if cfg.StoreInterval == unknownIntFieldValue {
		flag2.StringVarP(iFlag, "i", "i", "", "")
	}

	if cfg.FileStoragePath == unknownStringFieldValue {
		flag2.StringVarP(fFlag, "f", "f", unknownStringFieldValue, "")
	}

	if cfg.ConnectionDB == unknownStringFieldValue {
		flag2.StringVarP(dFlag, "d", "d", unknownStringFieldValue, "")
	}
}

func processServerFlags(cfg *ServerConfig) error {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	dFlag, keyFlag := unknownStringFieldValue, unknownStringFieldValue
	address, rFlag, fFlag := unknownStringFieldValue, unknownStringFieldValue, unknownStringFieldValue

	var iFlag string
	addServerFlags(cfg, &address, &rFlag, &iFlag, &fFlag, &dFlag, &keyFlag)

	flag2.Parse()

	setUnknownStrValue(&cfg.Address, address)
	setUnknownStrValue(&cfg.HashKey, keyFlag)
	setUnknownStrValue(&cfg.FileStoragePath, fFlag)
	setUnknownStrValue(&cfg.ConnectionDB, dFlag)

	if !cfg.hasRestore && rFlag != unknownStringFieldValue {
		if s, err := strconv.ParseBool(rFlag); err == nil {
			cfg.hasRestore = true
			cfg.Restore = s
		} else {
			return fmt.Errorf("couldn't convert 'r' to bool, rFlag: %s, err: %w", rFlag, err)
		}
	}

	if cfg.StoreInterval == unknownIntFieldValue && iFlag != "" {
		if s, err := strconv.ParseInt(iFlag, 10, 64); err == nil {
			cfg.StoreInterval = s
		} else {
			return fmt.Errorf("couldn't convert the store interval to int, iFlag: %s, err: %w", iFlag, err)
		}
	}

	return nil
}

func setUnknownStrValue(target *string, value string) {
	if value != unknownStringFieldValue {
		*target = value
	}
}

func addAgentFlags(cfg *AgentConfig, address *string,
	hashKey *string, pollInterval *string, reportInterval *string, limit *string,
) {
	if cfg.Address == "" {
		flag2.StringVarP(address, "a", "a", "", "")
	}

	if cfg.HashKey == unknownStringFieldValue {
		flag2.StringVarP(hashKey, "k", "k", "", "")
	}

	if cfg.PollInterval == 0 {
		flag2.StringVarP(pollInterval, "p", "p", "", "")
	}

	if cfg.ReportInterval == 0 {
		flag2.StringVarP(reportInterval, "r", "r", "", "")
	}

	if cfg.RateLimit == 0 {
		flag2.StringVarP(limit, "l", "l", "", "")
	}
}

func processAgentFlags(cfg *AgentConfig) error {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true

	var address, keyFlag, pFlag, rFlag, lFlag string

	addAgentFlags(cfg, &address, &keyFlag, &pFlag, &rFlag, &lFlag)
	flag2.Parse()

	if address != "" {
		cfg.Address = address
	}

	if keyFlag != "" {
		cfg.HashKey = keyFlag
	}

	if err := setAgentIntFlag(&cfg.PollInterval, pFlag, "poll interval"); err != nil {
		return err
	}
	if err := setAgentIntFlag(&cfg.ReportInterval, rFlag, "request interval"); err != nil {
		return err
	}
	err := setAgentIntFlag(&cfg.RateLimit, lFlag, "rate limit")

	return err
}

func setAgentIntFlag(cfgInt *int64, flag string, errMesPart string) error {
	if *cfgInt == 0 && flag != "" {
		if s, err := strconv.ParseInt(flag, 10, 64); err == nil {
			*cfgInt = s
		} else {
			return fmt.Errorf("couldn't convert the %s to int, flag: %s, err: %w", errMesPart, flag, err)
		}
	}

	return nil
}

func ProcessEnvServer(config *ServerConfig) error {
	log.Println(os.Environ())

	opts := env.Options{ //nolint:exhaustruct
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if tag == "RESTORE" {
				config.hasRestore = true
			}
			log.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}

	if err := env.Parse(config, opts); err != nil {
		return fmt.Errorf("failed to parse an environment, error: %w", err)
	}

	return nil
}

var processEnvAgent = func(config *AgentConfig) error {
	err := env.Parse(config)
	if err != nil {
		return fmt.Errorf("failed to parse an environment, error: %w", err)
	}

	return nil
}

func setupDefaultServerValues(config *ServerConfig,
	defaultAddress string,
	defaultStoreInterval int64,
	defaultFileStoragePath string,
	defaultKey string,
	defaultRestore bool,
) {
	if config.Address == unknownStringFieldValue {
		config.Address = defaultAddress
	}

	if config.HashKey == unknownStringFieldValue {
		config.HashKey = defaultKey
	}

	if config.StoreInterval == unknownIntFieldValue {
		config.StoreInterval = defaultStoreInterval
	}

	if config.FileStoragePath == unknownStringFieldValue {
		config.FileStoragePath = defaultFileStoragePath
	}

	if !config.hasRestore {
		config.Restore = defaultRestore
	}
}

func setupDefaultAgentValues(config *AgentConfig,
	defaultAddress string,
	defaultRepInterval time.Duration,
	defaultPollInterval time.Duration,
) {
	if config.HashKey == unknownStringFieldValue {
		config.HashKey = ""
	}

	if config.Address == "" {
		config.Address = defaultAddress
	}

	if config.ReportInterval == 0 {
		config.ReportInterval = int64(defaultRepInterval)
	}
	if config.PollInterval == 0 {
		config.PollInterval = int64(defaultPollInterval)
	}
}

// IsUseDatabase shows an ability to use a DB by ServerConfig.
func (cfg ServerConfig) IsUseDatabase() bool {
	return cfg.ConnectionDB != "" && cfg.ConnectionDB != unknownStringFieldValue
}

func (cfg ServerConfig) String() string {
	return fmt.Sprintf("Address: %s \n StoreInterval: %d \n"+
		" FileStoragePath: %s \n"+
		" ConnectionDB: %s \n"+
		" Key: %s \n"+
		" Restore: %t \n",
		cfg.Address, cfg.StoreInterval, cfg.FileStoragePath, cfg.ConnectionDB, cfg.HashKey, cfg.Restore)
}
