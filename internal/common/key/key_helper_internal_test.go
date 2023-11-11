package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wdTest = ""

func getWD() string {
	if wdTest != "" {
		return wdTest
	}
	wDir, _ := os.Getwd()
	//nolint:forbidigo
	fmt.Println("wd started:")
	for !strings.HasSuffix(wDir, "practicummetrics") {
		wDir = filepath.Dir(wDir)
	}
	wdTest = wDir

	return wDir
}

func TestLoadPrivateKey(t *testing.T) {
	wDir := getWD()
	filePath := fmt.Sprintf("%s/keys/keyfile.pem", wDir)
	key, err := loadPrivateKey(filePath)
	if !assert.NoError(t, err) {
		return
	}
	//nolint:forbidigo
	fmt.Println(key.N)
}

func TestLoadPublicKey(t *testing.T) {
	wDir := getWD()

	filePath := fmt.Sprintf("%s/keys/publickeyfile.pem", wDir)
	key, err := loadPublicKey(filePath)
	if !assert.NoError(t, err) {
		return
	}
	//nolint:forbidigo
	fmt.Println(key.N)
}

func TestEncryptDecrypt(t *testing.T) {
	wDir := getWD()

	keyPu, err := loadPublicKey(fmt.Sprintf("%s/keys/publickeyfile.pem", wDir))
	if !assert.NoError(t, err) {
		return
	}

	keyPr, err := loadPrivateKey(fmt.Sprintf("%s/keys/keyfile.pem", wDir))
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
