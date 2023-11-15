package repository

import (
	"crypto/aes"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/key"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	//nolint:exhaustruct
	cfgS := config.ServerConfig{
		Config: config.Config{
			CryptoKey: "keys/keyfile.pem",
		},
	}
	keyPr, err := key.LoadPrivateKey(cfgS)
	require.NoError(t, err)
	require.NotNil(t, keyPr)
	SetPrivateKey(keyPr)

	//nolint:exhaustruct
	cfgA := config.AgentConfig{
		Config: config.Config{
			CryptoKey: "keys/publickeyfile.pem",
		},
	}
	keyPub, err := key.LoadPublicKey(cfgA)
	require.NoError(t, err)
	require.NotNil(t, keyPub)
	SetPublicKey(keyPub)

	want := "test message"

	encB, err := encryptByPublicKey(keyPub, []byte(want))
	assert.NoError(t, err)

	decB, err := decryptByPrivateKey(encB)
	assert.NoError(t, err)
	assert.Equal(t, want, string(decB))
}

func TestEncryptBigMessageErr(t *testing.T) {
	_, err := EncryptBigMessage([]byte("dshua"))
	assert.Error(t, err)
}

func TestEncryptDecryptBigMessage(t *testing.T) {
	// Init

	//nolint:exhaustruct
	cfgS := config.ServerConfig{
		Config: config.Config{
			CryptoKey: "keys/keyfile.pem",
		},
	}
	LoadPrivateKey(cfgS)

	//nolint:exhaustruct
	cfgA := config.AgentConfig{
		Config: config.Config{
			CryptoKey: "keys/publickeyfile.pem",
		},
	}
	err := InitAgentAesKeys(cfgA)
	require.NoError(t, err)

	want := "a big message"

	mess, err := EncryptBigMessage([]byte(want))
	assert.NoError(t, err)
	got, err := DecryptBigMessage(mess.Encoded, mess.AesKey)
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func TestLoadPrivateKeyErr(t *testing.T) {
	keySt.private = nil
	//nolint:exhaustruct
	cfgS := config.ServerConfig{
		Config: config.Config{
			CryptoKey: "",
		},
	}
	LoadPrivateKey(cfgS)

	assert.Nil(t, keySt.private)
}

func TestInitAgentAesKeys(t *testing.T) {
	//nolint:exhaustruct
	cfga := config.AgentConfig{
		Config: config.Config{
			CryptoKey: "keys/publickeyfile.pem",
		},
	}
	err := InitAgentAesKeys(cfga)

	assert.NoError(t, err)
}

func TestInitAgentAesKeysErrorLoadRSA(t *testing.T) {
	//nolint:exhaustruct
	cfga := config.AgentConfig{
		Config: config.Config{
			CryptoKey: "",
		},
	}
	err := InitAgentAesKeys(cfga)

	assert.Error(t, err)
}

func TestAesDecodeBadKey(t *testing.T) {
	got, err := aesDecode(nil, []byte(""))

	assert.Nil(t, got)
	assert.ErrorIs(t, err, aes.KeySizeError(0))
}
