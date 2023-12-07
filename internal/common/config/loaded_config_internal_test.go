package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrimLastSOnlyS(t *testing.T) {
	income := "s"
	want := ""

	got := trimLastS(income)

	assert.Equal(t, want, got)
}

func TestTrimLastS(t *testing.T) {
	income := "q"
	want := "q"

	got := trimLastS(income)

	assert.Equal(t, want, got)
}

func TestLoadAgentConfigValues(t *testing.T) {
	loadedCfg := LoadedAgentConfig{
		Address:        "localhost:8888",
		ReportInterval: "3s",
		PollInterval:   "4s",
		CryptoKey:      "keys22/publickeyfile.pem",
	}

	//nolint:exhaustruct
	want := AgentConfig{
		Config:         Config{Address: "localhost:8888", CryptoKey: "keys22/publickeyfile.pem"},
		ReportInterval: 3,
		PollInterval:   4,
	}

	aCfg := new(AgentConfig)
	aCfg.Address = unknownStringFieldValue
	aCfg.CryptoKey = unknownStringFieldValue

	got := fillAgentConfigIfEmpty(*aCfg, loadedCfg)

	assert.Equal(t, want, got)
}

func TestLoadAgentConfigFile(t *testing.T) {
	want := LoadedAgentConfig{
		Address:        "localhost:8080",
		ReportInterval: "1s",
		PollInterval:   "1s",
		CryptoKey:      "keys/publickeyfile.pem",
	}
	wDir := common.GetWD()
	path := fmt.Sprintf("%s/testdata/agent_example_cfg.json", wDir)
	jsonString, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NotNil(t, jsonString)
	lacfg := LoadedAgentConfig{} //nolint:exhaustruct
	err = json.Unmarshal(jsonString, &lacfg)
	assert.NoError(t, err)
	assert.Equal(t, want, lacfg)
}

func TestLoadServerConfigFile(t *testing.T) {
	//nolint:exhaustruct
	want := LoadedServerConfig{
		Address:       "localhost:8080",
		Restore:       true,
		StoreInterval: "1s",
		StoreFile:     "/path/to/file.db",
		CryptoKey:     "keys/keyfile.pem",
	}
	wDir := common.GetWD()
	path := fmt.Sprintf("%s/testdata/server_example_cfg.json", wDir)
	jsonString, err := os.ReadFile(path)
	require.NoError(t, err)
	require.NotNil(t, jsonString)
	lacfg := LoadedServerConfig{} //nolint:exhaustruct
	err = json.Unmarshal(jsonString, &lacfg)
	assert.NoError(t, err)
	assert.Equal(t, want, lacfg)
}

func TestLoadServerConfigValues(t *testing.T) {
	loadedCfg := LoadedServerConfig{
		Address:       "localhost:8888",
		Restore:       true,
		StoreInterval: "22s",
		StoreFile:     "/tmp/store.txt",
		DatabaseDsn:   "nodbhere:889862",
		CryptoKey:     "keys22/keyfile.pem",
	}

	//nolint:exhaustruct
	want := ServerConfig{
		Config: Config{
			Address:    "localhost:8888",
			CryptoKey:  "keys22/keyfile.pem",
			ConfigFile: unknownStringFieldValue,
		},
		StoreInterval:   22,
		FileStoragePath: "/tmp/store.txt",
		ConnectionDB:    "nodbhere:889862",
		TrustedSubnet:   unknownStringFieldValue,
		Restore:         true,
		hasRestore:      true,
	}

	got := NewServerConfig()
	got.HashKey = ""

	fillServerConfigIfEmpty(got, loadedCfg)

	assert.Equal(t, want, *got)
}
