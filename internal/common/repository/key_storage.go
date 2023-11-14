package repository

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/key"
	"github.com/DimaKoz/practicummetrics/internal/common/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	keyStorageSync = &sync.RWMutex{}
	keySt          = keyStorage{} //nolint:exhaustruct
)

// keyStorage represents storage.
type keyStorage struct {
	private            *rsa.PrivateKey
	public             *rsa.PublicKey
	aesEncodeKey       []byte
	aesKeyEncodedByRsa string // there is a key prepared to send
}

// decryptByPrivateKey returns a decrypted []byte or nil with an error.
func decryptByPrivateKey(encryptedMessage []byte) ([]byte, error) {
	//nolint:exhaustruct
	cryptoOps := rsa.OAEPOptions{Hash: crypto.SHA256}

	decryptedBytes, err := keySt.private.Decrypt(nil, encryptedMessage, &cryptoOps)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt by: ")
	}

	return decryptedBytes, nil
}

// encryptByPublicKey returns an encrypted []byte or nil with an error.
func encryptByPublicKey(keyRsa *rsa.PublicKey, message []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		keyRsa,
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

func initAesKey(buffSize int) ([]byte, error) {
	data := make([]byte, buffSize)
	if _, err := rand.Read(data); err != nil {
		return nil, fmt.Errorf("can't create AES key by : %w", err)
	}

	return data, nil
}

func InitAgentAesKeys(cfg config.AgentConfig) error {
	const bufSizeAes128 = 16
	keyAes, err := initAesKey(bufSizeAes128 /*AES-128*/)
	if err != nil {
		return err
	}
	keyStorageSync.Lock()
	defer keyStorageSync.Unlock()
	keySt.aesEncodeKey = keyAes

	keyPublic, err := key.LoadPublicKey(cfg)
	if err != nil {
		return errors.Wrap(err, "can't load RSA public key: ")
	}
	keySt.public = keyPublic
	encryptedKey, err := encryptByPublicKey(keySt.public, keyAes)
	if err != nil {
		return errors.Wrap(err, "can't encrypt AES key: ")
	}
	keySt.aesKeyEncodedByRsa = base64.RawStdEncoding.EncodeToString(encryptedKey)

	return nil
}

func EncryptBigMessage(bigMessage []byte) (*model.EncMessage, error) {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()
	aesKey := keySt.aesEncodeKey

	encBigMessage, err := aesEncode(aesKey, bigMessage)
	if err != nil {
		return nil, errors.Wrap(err, "can't encrypt a big message by: ")
	}
	msg := make([]byte, base64.RawStdEncoding.EncodedLen(len(encBigMessage)))
	base64.RawStdEncoding.Encode(msg, encBigMessage)

	result := &model.EncMessage{Encoded: msg, AesKey: keySt.aesKeyEncodedByRsa}

	return result, nil
}

func DecryptBigMessage(bigMessage []byte, encodedAesKey string) ([]byte, error) {
	keyStorageSync.RLock()
	defer keyStorageSync.RUnlock()

	encryptedKey, err := base64.RawStdEncoding.DecodeString(encodedAesKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt a big message by: ")
	}

	encryptedMessage := make([]byte, base64.RawStdEncoding.DecodedLen(len(bigMessage)))
	_, err = base64.RawStdEncoding.Decode(encryptedMessage, bigMessage)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt a big message by: ")
	}

	aesKey, err := decryptByPrivateKey(encryptedKey)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt a big message by: ")
	}

	decryptedMessage, err := aesDecode(aesKey, encryptedMessage)
	if err != nil {
		return nil, errors.Wrap(err, "can't decrypt a big message by: ")
	}

	return decryptedMessage, nil
}

func aesDecode(aesKey []byte, message []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := message[:nonceSize], message[nonceSize:]

	decoded, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return decoded, nil
}

func aesEncode(aesKey []byte, message []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	ciphertext := gcm.Seal(nonce, nonce, message, nil)

	return ciphertext, nil
}
