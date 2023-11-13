package key

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/stretchr/testify/assert"
)

var testWantKeyN = "208841538677995080403984073650783025092403332474897996743858422571052822978339608048996" +
	"4004348626553067273762513574747900912286916180395909246475105889396959918521462496023938039204867605370" +
	"5871948538298442052892363849095055039887484117149661353849048085782503492826837547700799734398806478587" +
	"56120385701180030591969007272396670968092138618326959619156081278324693924793471001360189297029932336155" +
	"21228800225119530304254936673132081730981569231768762686704094429887729103363188504940360091508434836712" +
	"53030991843518194730978606046862655717821804120033175666724338968755482322571009579712940573244612415306055690777991"

func TestLoadPrivateKey1(t *testing.T) {
	//nolint:exhaustruct
	cfg := config.ServerConfig{
		Config: config.Config{
			CryptoKey: "keys/keyfile.pem",
		},
	}
	key, err := LoadPrivateKey(cfg)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, key)
	assert.Equal(t, testWantKeyN, key.N.String(), "key.N must be equal to the provided key")
}

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

func TestLoadPublicKey1(t *testing.T) {
	//nolint:exhaustruct
	cfg := config.AgentConfig{
		Config: config.Config{
			CryptoKey: "keys/publickeyfile.pem",
		},
	}
	key, err := LoadPublicKey(cfg)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, key)
	assert.Equal(t, testWantKeyN, key.N.String(), "key.N must be equal to the provided key")
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
