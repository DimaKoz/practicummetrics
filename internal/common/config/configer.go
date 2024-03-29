package config

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

// Constants for configs.
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
	Address    string `env:"ADDRESS"`
	HashKey    string `env:"KEY"`
	CryptoKey  string `env:"CRYPTO_KEY"`
	ConfigFile string `env:"CONFIG"`
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
	TrustedSubnet   string `env:"TRUSTED_SUBNET"`
	hasRestore      bool
	Restore         bool `env:"RESTORE"`
}

// NewServerConfig creates an instance of ServerConfig.
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Config: Config{
			Address:    unknownStringFieldValue,
			HashKey:    unknownStringFieldValue,
			CryptoKey:  unknownStringFieldValue,
			ConfigFile: unknownStringFieldValue,
		},
		StoreInterval:   unknownIntFieldValue,
		FileStoragePath: unknownStringFieldValue,
		ConnectionDB:    unknownStringFieldValue,
		TrustedSubnet:   unknownStringFieldValue,
		hasRestore:      false,
		Restore:         true,
	}
}

func (cfg ServerConfig) HasTrustedSubnet() bool {
	return cfg.TrustedSubnet != "" && cfg.TrustedSubnet != unknownStringFieldValue
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

	if cfg.ConfigFile != unknownStringFieldValue && cfg.ConfigFile != "" {
		laCfg, err := LoadConfigFromFile[LoadedServerConfig](cfg.ConfigFile)
		if err != nil {
			return fmt.Errorf("cannot load flags variables: %w", err)
		}
		fillServerConfigIfEmpty(cfg, *laCfg)
	}

	setupDefaultServerValues(cfg,
		defaultAddress,
		defaultStoreInterval,
		defaultFileStoragePath,
		defaultKey,
		defaultRestore)

	return nil
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{} //nolint:exhaustruct
	cfg.HashKey = unknownStringFieldValue
	cfg.CryptoKey = unknownStringFieldValue
	cfg.ConfigFile = unknownStringFieldValue

	return cfg
}

// LoadAgentConfig returns *AgentConfig.
func LoadAgentConfig() (*AgentConfig, error) {
	cfg := NewAgentConfig()

	if err := processEnvAgent(cfg); err != nil {
		return nil, fmt.Errorf("agent config: cannot process ENV variables: %w", err)
	}

	if err := processAgentFlags(cfg); err != nil {
		return nil, fmt.Errorf("cannot process flags variables: %w", err)
	}
	if cfg.ConfigFile != unknownStringFieldValue && cfg.ConfigFile != "" {
		laCfg, err := LoadConfigFromFile[LoadedAgentConfig](cfg.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("cannot load flags variables: %w", err)
		}
		cfgNew := fillAgentConfigIfEmpty(*cfg, *laCfg)
		cfg = &cfgNew
	}

	setupDefaultAgentValues(cfg, defaultAddress, defaultReportInterval, defaultPollInterval)

	return cfg, nil
}

// addServerFlags adds server flags to process them.
//
//nolint:cyclop
func addServerFlags(cfg *ServerConfig,
	address, rFlag, iFlag, fFlag, dFlag, keyFlag, sFlag, cFlag, tFlag *string,
) {
	if cfg.Address == unknownStringFieldValue {
		flag.StringVar(address, "a", unknownStringFieldValue, "")
	}

	if cfg.HashKey == unknownStringFieldValue && flag.Lookup("k") == nil {
		flag.StringVar(keyFlag, "k", unknownStringFieldValue, "")
	}
	if cfg.CryptoKey == unknownStringFieldValue && flag.Lookup("crypto-key") == nil {
		flag.StringVar(sFlag, "crypto-key", unknownStringFieldValue, "")
	}
	addStringChecksStringFlag(cfg.ConfigFile, unknownStringFieldValue, "c", cFlag)
	addStringChecksStringFlag(cfg.ConfigFile, unknownStringFieldValue, "config", cFlag)

	if !cfg.hasRestore && flag.Lookup("r") == nil {
		flag.StringVar(rFlag, "r", unknownStringFieldValue, "")
	}

	addIntChecksStringFlag(cfg.StoreInterval, unknownIntFieldValue, "i", iFlag)

	if cfg.FileStoragePath == unknownStringFieldValue {
		flag.StringVar(fFlag, "f", unknownStringFieldValue, "")
	}

	if cfg.ConnectionDB == unknownStringFieldValue {
		flag.StringVar(dFlag, "d", unknownStringFieldValue, "")
	}

	if cfg.TrustedSubnet == unknownStringFieldValue {
		flag.StringVar(tFlag, "t", unknownStringFieldValue, "")
	}
}

// processServerFlags gets parameters from command line and fill ServerConfig
// or returns error if something wrong.
func processServerFlags(cfg *ServerConfig) error {
	dFlag, keyFlag := unknownStringFieldValue, unknownStringFieldValue
	address, rFlag, fFlag := unknownStringFieldValue, unknownStringFieldValue, unknownStringFieldValue
	sFlag, cFlag := unknownStringFieldValue, unknownStringFieldValue
	tFlag := unknownStringFieldValue

	var iFlag string
	addServerFlags(cfg, &address, &rFlag, &iFlag, &fFlag, &dFlag, &keyFlag, &sFlag, &cFlag, &tFlag)
	flag.Parse()

	setUnknownStrValue(&cfg.Address, address)
	setUnknownStrValue(&cfg.HashKey, keyFlag)
	setUnknownStrValue(&cfg.FileStoragePath, fFlag)
	setUnknownStrValue(&cfg.ConnectionDB, dFlag)
	setUnknownStrValue(&cfg.CryptoKey, sFlag)
	setUnknownStrValue(&cfg.ConfigFile, cFlag)
	setUnknownStrValue(&cfg.TrustedSubnet, tFlag)

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

// setUnknownStrValue sets a value to a target.
func setUnknownStrValue(target *string, value string) {
	if value != unknownStringFieldValue {
		*target = value
	}
}

// addAgentFlags adds agent flags to process them.
func addAgentFlags(cfg *AgentConfig, address, hashKey, pollInterval, reportInterval, limit, cryptoKey, cFlag *string,
) {
	addStringChecksStringFlag(cfg.Address, "", "a", address)

	addStringChecksStringFlag(cfg.HashKey, unknownStringFieldValue, "k", hashKey)

	addStringChecksStringFlag(cfg.CryptoKey, unknownStringFieldValue, "crypto-key", cryptoKey)
	addStringChecksStringFlag(cfg.ConfigFile, unknownStringFieldValue, "c", cFlag)
	addStringChecksStringFlag(cfg.ConfigFile, unknownStringFieldValue, "config", cFlag)

	addIntChecksStringFlag(cfg.PollInterval, 0, "p", pollInterval)

	addIntChecksStringFlag(cfg.ReportInterval, 0, "r", reportInterval)

	addIntChecksStringFlag(cfg.RateLimit, 0, "l", limit)
}

func addStringChecksStringFlag(currentCfgValue, defaultCfgValue, flagName string, passedVar *string) {
	if currentCfgValue == defaultCfgValue && flag.Lookup(flagName) == nil {
		flag.StringVar(passedVar, flagName, "", "")
	}
}

func addIntChecksStringFlag(currentCfgValue, defaultCfgValue int64, flagName string, passedVar *string) {
	if currentCfgValue == defaultCfgValue && flag.Lookup(flagName) == nil {
		flag.StringVar(passedVar, flagName, "", "")
	}
}

// processAgentFlags gets parameters from command line and fill AgentConfig
// or returns error if something wrong.
func processAgentFlags(cfg *AgentConfig) error {
	flag.CommandLine.ErrorHandling()
	var address, keyFlag, pFlag, rFlag, lFlag, sFlag, cFlag string

	addAgentFlags(cfg, &address, &keyFlag, &pFlag, &rFlag, &lFlag, &sFlag, &cFlag)
	flag.Parse()

	if address != "" {
		cfg.Address = address
	}

	if keyFlag != "" {
		cfg.HashKey = keyFlag
	}

	if sFlag != "" {
		cfg.CryptoKey = sFlag
	}
	if cFlag != "" {
		cfg.ConfigFile = cFlag
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

// setAgentIntFlag sets flag value to int64 field of AgentConfig.
func setAgentIntFlag(cfgInt *int64, flag, errMesPart string) error {
	if *cfgInt == 0 && flag != "" {
		if s, err := strconv.ParseInt(flag, 10, 64); err == nil {
			*cfgInt = s
		} else {
			return fmt.Errorf("couldn't convert the %s to int, flag: %s, err: %w", errMesPart, flag, err)
		}
	}

	return nil
}

// ProcessEnvServer fills ServerConfig.
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

// processEnvAgent an implementation of a function to parse AgentConfig.
var processEnvAgent = func(config *AgentConfig) error {
	err := env.Parse(config)
	if err != nil {
		return fmt.Errorf("failed to parse an environment, error: %w", err)
	}

	return nil
}

// setupDefaultServerValues sets default values to ServerConfig.
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

	if config.CryptoKey == unknownStringFieldValue {
		config.CryptoKey = defaultKey
	}

	if config.TrustedSubnet == unknownStringFieldValue {
		config.TrustedSubnet = ""
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

// setupDefaultAgentValues sets default values to AgentConfig.
func setupDefaultAgentValues(config *AgentConfig,
	defaultAddress string,
	defaultRepInterval time.Duration,
	defaultPollInterval time.Duration,
) {
	if config.HashKey == unknownStringFieldValue {
		config.HashKey = ""
	}
	if config.CryptoKey == unknownStringFieldValue {
		config.CryptoKey = ""
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

// String is Stringer implementation of ServerConfig.
func (cfg ServerConfig) String() string {
	return cfg.StringVariantCopy()
}

// StringVariantCopy a variant of string representation of ServerConfig.
func (cfg ServerConfig) StringVariantCopy() string {
	const minimumLen = 101
	storeI := strconv.FormatInt(cfg.StoreInterval, 10)
	restore := strconv.FormatBool(cfg.Restore)
	grow := minimumLen +
		len(storeI) + len(restore) + len(cfg.Address) + len(cfg.FileStoragePath) + len(cfg.ConnectionDB) + len(cfg.HashKey)
	result := make([]byte, grow)
	bLen := 0
	bLen += copy(result[bLen:], "Address: ")
	bLen += copy(result[bLen:], cfg.Address)
	bLen += copy(result[bLen:], " \n StoreInterval: ")
	bLen += copy(result[bLen:], storeI)
	bLen += copy(result[bLen:], " \n FileStoragePath: ")
	bLen += copy(result[bLen:], cfg.FileStoragePath)
	bLen += copy(result[bLen:], " \n ConnectionDB: ")
	bLen += copy(result[bLen:], cfg.ConnectionDB)
	bLen += copy(result[bLen:], " \n Key: ")
	bLen += copy(result[bLen:], cfg.HashKey)
	bLen += copy(result[bLen:], " \n CryptoKey: ")
	bLen += copy(result[bLen:], cfg.CryptoKey)
	bLen += copy(result[bLen:], " \n Restore: ")
	bLen += copy(result[bLen:], restore)
	_ = copy(result[bLen:], " \n")

	return string(result)
}

func PrepBuildValues(bldV, bldD, bldC string) string {
	buffStr := bytes.Buffer{}
	buffStr.WriteString("Build version: ")
	buffStr.WriteString(bldV)
	buffStr.WriteString("\nBuild date: ")
	buffStr.WriteString(bldD)
	buffStr.WriteString("\nBuild commit:  ")
	buffStr.WriteString(bldC)
	buffStr.WriteString("\n")

	return buffStr.String()
}
