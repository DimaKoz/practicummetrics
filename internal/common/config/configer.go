package config

import (
	"fmt"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/caarlos0/env/v6"
	flag2 "github.com/spf13/pflag"
	"strconv"
	"time"
)

func AgentInitConfig(cfg *model.Config,
	defaultAddress string,
	defaultRepInterval time.Duration,
	defaultPollInterval time.Duration) {

	if cfg == nil {
		fmt.Println("passed config is nil")
		return
	}

	processEnv(cfg)

	processFlags(cfg)

	setupDefaultValues(cfg, defaultAddress, defaultRepInterval, defaultPollInterval)
}

func processFlags(cfg *model.Config) {
	flag2.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var address string
	if cfg.Address == "" {
		flag2.StringVarP(&address, "a", "a", "", "")
		cfg.Address = address
	}

	var pFlag string
	if cfg.PollInterval == 0 {
		flag2.StringVarP(&pFlag, "p", "p", "", "")
	}

	var rFlag string
	if cfg.PollInterval == 0 {
		flag2.StringVarP(&rFlag, "r", "r", "", "")
	}

	flag2.Parse()
	if address != "" {
		cfg.Address = address
	}

	if s, err := strconv.ParseInt(pFlag, 10, 64); err == nil {
		cfg.PollInterval = s
	} else {
		fmt.Println("pFlag:", pFlag, ", s:", s, ", err:", err)
	}

	if s, err := strconv.ParseInt(rFlag, 10, 64); err == nil {
		cfg.ReportInterval = s
	} else {
		fmt.Println("rFlag:", rFlag, ", s:", s, ", err:", err)
	}
}

func processEnv(config *model.Config) {
	err := env.Parse(config)
	if err != nil {
		fmt.Println(" env parsing error: ", err)
	}
}

func setupDefaultValues(config *model.Config,
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
