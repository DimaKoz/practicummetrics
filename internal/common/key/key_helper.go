package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var errWrongKeyType = errors.New("wrong key type")

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key, ok := parseResult.(*rsa.PrivateKey)
	if !ok {
		return nil, errWrongKeyType
	}

	return key, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKIXPublicKey(block.Bytes)
	key, ok := parseResult.(*rsa.PublicKey)
	if !ok {
		return nil, errWrongKeyType
	}

	return key, nil
}
