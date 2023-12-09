package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/pkg/errors"
)

var ErrNoKeyPath = errors.New("no crypto-key path")

func LoadPrivateKey(cfg config.ServerConfig) (*rsa.PrivateKey, error) {
	if cfg.CryptoKey == "" {
		return nil, ErrNoKeyPath
	}
	filePath := fmt.Sprintf("%s/%s", common.GetWD(), cfg.CryptoKey /*"keys/keyfile.pem"*/)

	return loadPrivateKeyImpl(filePath)
}

const errLoadingKeyTemplate = "can't load a key by: %w"

func loadPrivateKeyImpl(path string) (*rsa.PrivateKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(errLoadingKeyTemplate, err)
	}
	block, _ := pem.Decode(pemString)
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf(errLoadingKeyTemplate, err)
	}
	key, _ := parseResult.(*rsa.PrivateKey)

	return key, nil
}

func LoadPublicKey(cfg config.AgentConfig) (*rsa.PublicKey, error) {
	if cfg.CryptoKey == "" {
		return nil, ErrNoKeyPath
	}
	filePath := fmt.Sprintf("%s/%s", common.GetWD(), cfg.CryptoKey /*"keys/publickeyfile.pem"*/)

	return loadPublicKeyImpl(filePath)
}

func loadPublicKeyImpl(path string) (*rsa.PublicKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKIXPublicKey(block.Bytes)
	key, _ := parseResult.(*rsa.PublicKey)

	return key, nil
}
