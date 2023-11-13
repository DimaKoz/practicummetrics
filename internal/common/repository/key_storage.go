package repository

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"sync"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/key"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	keyStorageSync = &sync.RWMutex{}
	keySt          = keyStorage{} //nolint:exhaustruct
)

// keyStorage represents storage.
type keyStorage struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

// DecryptByPrivateKey returns a decrypted []byte or nil with an error.
func DecryptByPrivateKey(encryptedMessage []byte) ([]byte, error) {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()
	//nolint:exhaustruct
	cryptoOps := rsa.OAEPOptions{Hash: crypto.SHA256}

	decryptedBytes, err := keySt.private.Decrypt(nil, encryptedMessage, &cryptoOps)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt by: ")
	}

	return decryptedBytes, nil
}

// EncryptByPublicKey returns an encrypted []byte or nil with an error.
func EncryptByPublicKey(message []byte) ([]byte, error) {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		keySt.public,
		message,
		nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't encrypt by: ")
	}

	return encryptedBytes, nil
}

// SetPrivateKey sets rsa.PrivateKey.
func SetPrivateKey(key *rsa.PrivateKey) {
	keyStorageSync.Lock()
	defer keyStorageSync.Unlock()

	keySt.private = key
}

// SetPublicKey sets rsa.PublicKey.
func SetPublicKey(key *rsa.PublicKey) {
	keyStorageSync.Lock()
	defer keyStorageSync.Unlock()

	keySt.public = key
}

func LoadPrivateKey(cfg config.ServerConfig) {
	keyPrivate, err := key.LoadPrivateKey(cfg)
	if err != nil {
		zap.S().Error(err)

		return
	}
	SetPrivateKey(keyPrivate)
}

func LoadPublicKey(cfg config.AgentConfig) {
	keyPublic, err := key.LoadPublicKey(cfg)
	if err != nil {
		zap.S().Error(err)

		return
	}
	SetPublicKey(keyPublic)
}
