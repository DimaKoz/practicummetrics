package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPrivateKey(t *testing.T) {
	wDir := getWD()
	filePath := fmt.Sprintf("%s/keys/keyfile.pem", wDir)
	key, err := loadPrivateKeyImpl(filePath)
	if !assert.NoError(t, err) {
		return
	}
	//nolint:forbidigo
	fmt.Println(key.N)
}

func TestLoadPublicKey(t *testing.T) {
	wDir := getWD()

	filePath := fmt.Sprintf("%s/keys/publickeyfile.pem", wDir)
	key, err := loadPublicKeyImpl(filePath)
	if !assert.NoError(t, err) {
		return
	}
	//nolint:forbidigo
	fmt.Println(key.N)
}

func TestEncryptDecrypt(t *testing.T) {
	wDir := getWD()

	keyPu, err := loadPublicKeyImpl(fmt.Sprintf("%s/keys/publickeyfile.pem", wDir))
	if !assert.NoError(t, err) {
		return
	}

	keyPr, err := loadPrivateKeyImpl(fmt.Sprintf("%s/keys/keyfile.pem", wDir))
	if !assert.NoError(t, err) {
		return
	}

	want := "super secret message"

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		keyPu,
		[]byte(want),
		nil)
	if !assert.NoError(t, err) {
		return
	}
	//nolint:exhaustruct
	cryptoOps := rsa.OAEPOptions{Hash: crypto.SHA256}

	decryptedBytes, err := keyPr.Decrypt(nil, encryptedBytes, &cryptoOps)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, want, string(decryptedBytes))
}

func TestFiledLoadKeys(t *testing.T) {
	filePath := "badpath"
	_, err := loadPrivateKeyImpl(filePath)
	if !assert.Error(t, err) {
		return
	}
	_, err = loadPublicKeyImpl(filePath)
	if !assert.Error(t, err) {
		return
	}
}
