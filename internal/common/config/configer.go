package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	flag2 "github.com/spf13/pflag"
	"strconv"
	"time"
)

const (
	defaultPollInterval   = time.Duration(2)
	defaultReportInterval = time.Duration(10)
	defaultAddress        = "localhost:8080"

	defaultStoreInterval      = 300
	synchronicalStoreInterval = 0
	unknownIntFieldValue      = -1
	defaultFileStoragePath    = "/tmp/metrics-db.json"
	unknownStringFieldValue   = "unknownStringFieldValue"
	defaultRestore            = true
)

// Config represents a config of the agent and/or the server
type Config struct {
	Address string `env:"ADDRESS"`
}

type AgentConfig struct {
	Config
	ReportInterval int64 `env:"REPORT_INTERVAL"`
	PollInterval   int64 `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Config
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	hasRestore      bool
	Restore         bool `env:"RESTORE"`
}

func (cfg *ServerConfig) setupInitialServer() {
	cfg.Address = defaultAddress
	cfg.StoreInterval = defaultStoreInterval
	cfg.FileStoragePath = defaultFileStoragePath
	cfg.Restore = defaultRestore
}

func LoadServerConfig() (*ServerConfig, error) {

	cfg := &ServerConfig{
		Config:          Config{Address: unknownStringFieldValue},
		StoreInterval:   unknownIntFieldValue,
		FileStoragePath: unknownStringFieldValue,
		hasRestore:      false,
		Restore:         true,
	}

	if err := processEnv(cfg); err != nil {
		return nil, fmt.Errorf("server config: cannot process ENV variables: %w", err)
	}
	if err := processServerFlags(cfg); err != nil {
		return nil, fmt.Errorf("server config: cannot process flags variables: %w", err)
	}
	setupDefaultServerValues(cfg,
		defaultAddress,
		defaultStoreInterval,
		defaultFileStoragePath,
		defaultRestore)

	return cfg, nil
}

func LoadAgentConfig() (*AgentConfig, error) {
	cfg := &AgentConfig{}

	if err := processEnvAgent(cfg); err != nil {
		return nil, fmt.Errorf("agent config: cannot process ENV variables: %w", err)
	}
	if err := processAgentFlags(cfg); err != nil {
		return nil, fmt.Errorf("cannot process flags variables: %w", err)
	}
	setupDefaultAgentValues(cfg, defaultAddress, defaultReportInterval, defaultPollInterval)
	return cfg, nil
}

func processServerFlags(cfg *ServerConfig) error {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var address string
	if cfg.Address == unknownStringFieldValue {
		flag2.StringVarP(&address, "a", "a", unknownStringFieldValue, "")
	}

	var rFlag string
	if !cfg.hasRestore {
		flag2.StringVarP(&rFlag, "r", "r", unknownStringFieldValue, "")
	}

	var iFlag string
	if cfg.StoreInterval == unknownIntFieldValue {
		flag2.StringVarP(&iFlag, "i", "i", "", "")
	}

	var fFlag string
	if cfg.FileStoragePath == unknownStringFieldValue {
		flag2.StringVarP(&fFlag, "f", "f", "unknownStringFieldValue", "")
	}

	flag2.Parse()

	if address != unknownStringFieldValue {
		cfg.Address = address
	}

	if fFlag != unknownStringFieldValue {
		cfg.FileStoragePath = fFlag
	}

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

func processAgentFlags(cfg *AgentConfig) error {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var address string
	var pFlag string
	var rFlag string

	if cfg.Address == "" {
		flag2.StringVarP(&address, "a", "a", "", "")
	}

	if cfg.PollInterval == 0 {
		flag2.StringVarP(&pFlag, "p", "p", "", "")
	}

	if cfg.ReportInterval == 0 {
		flag2.StringVarP(&rFlag, "r", "r", "", "")
	}

	flag2.Parse()

	if address != "" {
		cfg.Address = address
	}

	if cfg.PollInterval == 0 && pFlag != "" {
		if s, err := strconv.ParseInt(pFlag, 10, 64); err == nil {
			cfg.PollInterval = s
		} else {
			return fmt.Errorf("couldn't convert the poll interval to int, pFlag: %s, err: %w", pFlag, err)
		}
	}

	if cfg.ReportInterval == 0 && rFlag != "" {
		if s, err := strconv.ParseInt(rFlag, 10, 64); err == nil {
			cfg.ReportInterval = s
		} else {
			return fmt.Errorf("couldn't convert the request interval to int, rFlag: %s, err: %w", rFlag, err)
		}
	}

	return nil
}

var processEnv = func(config *ServerConfig) error {
	opts := env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if tag == "RESTORE" {
				config.hasRestore = true
			}
			//log.Printf("Set %s to %v (default? %v)\n", tag, value, isDefault)
		},
	}

	err := env.Parse(config, opts)
	if err != nil {
		return fmt.Errorf("couldn't parse an enviroment, error: %w", err)
	}
	return nil
}

var processEnvAgent = func(config *AgentConfig) error {
	err := env.Parse(config)
	if err != nil {
		return fmt.Errorf("couldn't parse an enviroment, error: %w", err)
	}
	return nil
}

func setupDefaultServerValues(config *ServerConfig,
	defaultAddress string,
	defaultStoreInterval int64,
	defaultFileStoragePath string,
	defaultRestore bool) {
	if config.Address == unknownStringFieldValue {
		config.Address = defaultAddress
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
	defaultPollInterval time.Duration) {
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

func (cfg ServerConfig) String() string {
	return "Address:" + cfg.Address + "\n" +
		"StoreInterval:" + strconv.FormatInt(cfg.StoreInterval, 10) + "\n" +
		"FileStoragePath:" + cfg.FileStoragePath + "\n" +
		"Restore:" + strconv.FormatBool(cfg.Restore) + "\n"
}
