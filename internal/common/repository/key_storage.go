package repository

import (
	"crypto/rsa"
	"sync"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/key"
	"go.uber.org/zap"
)

var (
	keyStorageSync = &sync.RWMutex{}
	keySt          = keyStorage{} //nolint:exhaustruct
)

// keyStorage represents storage.
type keyStorage struct {
	private *rsa.PrivateKey
	// public  *rsa.PublicKey
}

// SetPrivateKey sets rsa.PrivateKey.
func SetPrivateKey(key *rsa.PrivateKey) {
	keyStorageSync.Lock()
	defer keyStorageSync.Unlock()

	keySt.private = key
}

/*
// GetPrivateKey returns a *rsa.PrivateKey or nil.
func GetPrivateKey() *rsa.PrivateKey {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()

	return keySt.private
}
*/

/*
// SetPublicKey sets rsa.PublicKey.
func SetPublicKey(key *rsa.PublicKey) {
	keyStorageSync.Lock()
	defer keyStorageSync.Unlock()

	keySt.public = key
}
*/

/*
// GetPublicKey returns a *rsa.PublicKey or nil.
func GetPublicKey() *rsa.PublicKey {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()

	return keySt.public
}
*/

func LoadPrivateKey(cfg config.ServerConfig) {
	keyPrivate, err := key.LoadPrivateKey(cfg)
	if err != nil {
		zap.S().Error(err)

		return
	}
	SetPrivateKey(keyPrivate)
}

/*func LoadPublicKey(cfg config.AgentConfig) {
	if cfg.CryptoKey == "" {
		return
	}
}
*/
