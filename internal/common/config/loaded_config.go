package config

import (
	"strconv"
	"unicode/utf8"
)

type LoadedAgentConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"` //nolint:tagliatelle
	PollInterval   string `json:"poll_interval"`   //nolint:tagliatelle
	CryptoKey      string `json:"crypto_key"`      //nolint:tagliatelle
}

func trimLastS(income string) string {
	lastRune, lastRuneSize := utf8.DecodeLastRuneInString(income)
	if lastRune != 's' {
		return income
	}

	return income[:len(income)-lastRuneSize]
}

func fillConfigIfEmpty(cfg AgentConfig, loadedCfg LoadedAgentConfig) AgentConfig {
	if cfg.Address == unknownStringFieldValue && loadedCfg.Address != "" {
		cfg.Address = loadedCfg.Address
	}

	if cfg.CryptoKey == unknownStringFieldValue && loadedCfg.CryptoKey != "" {
		cfg.CryptoKey = loadedCfg.CryptoKey
	}

	fillConfigIfEmptyInt(&cfg, loadedCfg)

	return cfg
}

func fillConfigIfEmptyInt(cfg *AgentConfig, loadedCfg LoadedAgentConfig) {
	if cfg.PollInterval == 0 && loadedCfg.PollInterval != "" {
		prepP := trimLastS(loadedCfg.PollInterval)
		if pollInterval, err := strconv.ParseInt(prepP, 10, 64); err == nil {
			cfg.PollInterval = pollInterval
		}
	}

	prepR := trimLastS(loadedCfg.ReportInterval)
	if cfg.ReportInterval == 0 && prepR != "" {
		if repInterval, err := strconv.ParseInt(prepR, 10, 64); err == nil {
			cfg.ReportInterval = repInterval
		}
	}
}
