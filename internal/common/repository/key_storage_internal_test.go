package repository

import (
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

	encB, err := EncryptByPublicKey([]byte(want))
	assert.NoError(t, err)

	decB, err := DecryptByPrivateKey(encB)
	assert.NoError(t, err)
	assert.Equal(t, want, string(decB))
}
