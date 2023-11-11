package key

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	key, _ := parseResult.(*rsa.PrivateKey)

	return key, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	pemString, err := os.ReadFile(path)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	block, _ := pem.Decode(pemString)
	parseResult, _ := x509.ParsePKIXPublicKey(block.Bytes)
	key, _ := parseResult.(*rsa.PublicKey)

	return key, nil
}
