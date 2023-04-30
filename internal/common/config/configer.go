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
}

func LoadServerConfig() (*ServerConfig, error) {

	cfg := &ServerConfig{}

	if err := processEnv(cfg); err != nil {
		return nil, fmt.Errorf("server config: cannot process ENV variables: %w", err)
	}
	processServerFlags(cfg)

	if cfg.Address == "" {
		cfg.Address = defaultAddress
	}

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

func processServerFlags(cfg *ServerConfig) {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var address string
	if cfg.Address == "" {
		flag2.StringVarP(&address, "a", "a", "", "")
	}

	flag2.Parse()

	if address != "" {
		cfg.Address = address
	}
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
	err := env.Parse(config)
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
