package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	flag2 "github.com/spf13/pflag"
	"strconv"
	"time"
)

const (
	// ServerCfg used to pass it to CreateConfig for getting a config for the server
	ServerCfg = iota
	// AgentCfg used to pass it to CreateConfig for getting a config for the server
	AgentCfg
)

const (
	defaultPollInterval   = time.Duration(2)
	defaultReportInterval = time.Duration(10)
	defaultAddress        = "localhost:8080"
)

// Config represents a config of the agent and/or the server
type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

var getEnv = processEnv

func CreateConfig(configType int) (*Config, error) {
	if configType != ServerCfg && configType != AgentCfg {
		errDesc := fmt.Sprintln("from CreateConfig: an unknown type of the config, no config for you ")
		return nil, fmt.Errorf(errDesc)
	}
	cfg := &Config{}

	warningEnv := getEnv(cfg)
	if warningEnv != nil {
		fmt.Println("from CreateConfig: warning: couldn't process the environment: ", warningEnv)
		fmt.Println("from CreateConfig: trying to use an another option")
	}
	fmt.Println("after getEnv:", cfg)
	warningFlags := processFlags(cfg, configType)
	if warningFlags != nil {
		fmt.Println("from CreateConfig: warning: couldn't process the flags: ", warningEnv)
		fmt.Println("from CreateConfig: trying to use default values")
	}
	setupDefaultValues(cfg, defaultAddress, defaultReportInterval, defaultPollInterval)

	return cfg, nil
}

func processFlags(cfg *Config, configType int) error {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var address string
	if cfg.Address == "" {
		flag2.StringVarP(&address, "a", "a", "", "")
		cfg.Address = address
	}
	var pFlag string
	var rFlag string
	if configType == AgentCfg {

		if cfg.PollInterval == 0 {
			flag2.StringVarP(&pFlag, "p", "p", "", "")
		}

		if cfg.PollInterval == 0 {
			flag2.StringVarP(&rFlag, "r", "r", "", "")
		}
	}

	flag2.Parse()

	if address != "" {
		cfg.Address = address
	}

	if configType == AgentCfg {

		if s, err := strconv.ParseInt(pFlag, 10, 64); err == nil {
			cfg.PollInterval = s
		} else {
			return fmt.Errorf("pFlag: %s, err: %w", pFlag, err)
		}

		if s, err := strconv.ParseInt(rFlag, 10, 64); err == nil {
			cfg.ReportInterval = s
		} else {
			return fmt.Errorf("rFlag: %s, err: %w", rFlag, err)
		}

	}
	return nil
}

func processEnv(config *Config) error {
	err := env.Parse(config)
	if err != nil {
		return fmt.Errorf(" env parsing error: %w", err)
	}
	return nil
}

func setupDefaultValues(config *Config,
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
