package config

import (
	"strconv"
	"unicode/utf8"
)

type LoadedServerConfig struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval string `json:"store_interval"` //nolint:tagliatelle
	StoreFile     string `json:"store_file"`     //nolint:tagliatelle
	DatabaseDsn   string `json:"database_dsn"`   //nolint:tagliatelle
	CryptoKey     string `json:"crypto_key"`     //nolint:tagliatelle
}

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

func fillAgentConfigIfEmpty(cfg AgentConfig, loadedCfg LoadedAgentConfig) AgentConfig {
	if cfg.Address == unknownStringFieldValue && loadedCfg.Address != "" {
		cfg.Address = loadedCfg.Address
	}

	if cfg.CryptoKey == unknownStringFieldValue && loadedCfg.CryptoKey != "" {
		cfg.CryptoKey = loadedCfg.CryptoKey
	}

	fillAgentConfigIfEmptyInt(&cfg, loadedCfg)

	return cfg
}

func fillAgentConfigIfEmptyInt(cfg *AgentConfig, loadedCfg LoadedAgentConfig) {
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

func fillServerConfigIfEmpty(cfg ServerConfig, loadedCfg LoadedServerConfig) ServerConfig {
	setUnknownStrValue(&cfg.Address, loadedCfg.Address)
	if loadedCfg.Address != "" {
		setUnknownStrValue(&cfg.Address, loadedCfg.Address)
	}

	if loadedCfg.CryptoKey != "" {
		setUnknownStrValue(&cfg.CryptoKey, loadedCfg.CryptoKey)
	}

	if loadedCfg.StoreFile != "" {
		setUnknownStrValue(&cfg.FileStoragePath, loadedCfg.StoreFile)
	}

	if loadedCfg.DatabaseDsn != "" {
		setUnknownStrValue(&cfg.ConnectionDB, loadedCfg.DatabaseDsn)
	}

	if cfg.StoreInterval == unknownIntFieldValue && loadedCfg.StoreInterval != "" {
		prepP := trimLastS(loadedCfg.StoreInterval)
		if storeInterval, err := strconv.ParseInt(prepP, 10, 64); err == nil {
			cfg.StoreInterval = storeInterval
		}
	}

	if !cfg.hasRestore && loadedCfg.Restore {
		cfg.Restore = loadedCfg.Restore
		cfg.hasRestore = true
	}

	return cfg
}
