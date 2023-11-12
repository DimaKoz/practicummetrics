package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var workDir = ""

func getWD() string {
	if workDir != "" {
		return workDir
	}
	wDirTemp, _ := os.Getwd()

	zap.S().Infof("wd started: %s", wDirTemp)

	for !strings.HasSuffix(wDirTemp, "practicummetrics") {
		wDirTemp = filepath.Dir(wDirTemp)
	}
	workDir = wDirTemp

	return workDir
}

var ErrNoKeyPath = errors.New("no crypto-key path")

func LoadPrivateKey(cfg config.ServerConfig) (*rsa.PrivateKey, error) {
	if cfg.CryptoKey == "" {
		return nil, ErrNoKeyPath
	}
	filePath := fmt.Sprintf("%s/%s", getWD(), cfg.CryptoKey /*"keys/keyfile.pem"*/)

	return loadPrivateKeyImpl(filePath)
}

func loadPrivateKeyImpl(path string) (*rsa.PrivateKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key, _ := parseResult.(*rsa.PrivateKey)

	return key, nil
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
