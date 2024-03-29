package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsBadHash(t *testing.T) {
	type args struct {
		cfgKey     string
		incomeHash string
		reqBody    []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "wrong",
			args: args{
				cfgKey:     "sdfs",
				incomeHash: "wrongHash",
				reqBody:    []byte("ABCDF"),
			},
			want: true,
		},
		{
			name: "ok",
			args: args{
				cfgKey:     "sdfs",
				incomeHash: "4f21bcb3b22ef261c261d033af3a8ad1fe8651f7edcc31180ac86a45b9040ee3",
				reqBody:    []byte("ABCDF"),
			},
			want: false,
		},
	}
	for _, unit := range tests {
		tt := unit
		t.Run(tt.name, func(t *testing.T) {
			got := isBadHash(tt.args.cfgKey, tt.args.incomeHash, tt.args.reqBody)
			assert.Equalf(t, tt.want, got, "isBadHash(%v, %v, %v)", tt.args.cfgKey, tt.args.incomeHash, tt.args.reqBody)
		})
	}
}

func TestGetRequestBody(t *testing.T) {
	const want = "abc"
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte(want)))
	eCtx := echoFr.AcquireContext()
	eCtx.SetRequest(request)
	got, gotErr := getRequestBody(eCtx)
	assert.NoError(t, gotErr)
	assert.Equal(t, want, string(got))
	err := echoFr.Close()
	require.NoError(t, err)
}

func TestAuthValidator(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	//nolint:exhaustruct
	cfg := config.ServerConfig{}
	comeIn := false
	authV := AuthValidator(cfg)(func(c echo.Context) error {
		comeIn = true

		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := authV(ctx)
	require.NoError(t, err)
	assert.True(t, comeIn)
}

func TestAuthValidatorErr(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	//nolint:exhaustruct
	cfg := config.ServerConfig{}
	comeIn := false
	authV := AuthValidator(cfg)(func(c echo.Context) error {
		comeIn = true

		return http.ErrAbortHandler
	})
	err := authV(ctx)
	assert.ErrorIs(t, err, http.ErrAbortHandler)
	assert.True(t, comeIn)
}

func TestAuthValidatorErrBadHash(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	//nolint:exhaustruct
	cfg := config.ServerConfig{Config: config.Config{HashKey: "gsgd"}}
	comeIn := false
	authV := AuthValidator(cfg)(func(c echo.Context) error {
		comeIn = true

		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := authV(ctx)
	assert.ErrorIs(t, err, errBadHash)
	assert.False(t, comeIn)
}

func TestAuthValidateEmptyCfgHashKey(t *testing.T) {
	echoFr := echo.New()
	eCtx := echoFr.AcquireContext()
	err := authValidate(eCtx, "")
	assert.NoError(t, err)
	err = echoFr.Close()
	require.NoError(t, err)
}

func TestAuthValidateEmptyHeaderHashKey(t *testing.T) {
	want := "abcd"
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte(want)))
	eCtx := echoFr.AcquireContext()
	eCtx.SetRequest(request)
	err := authValidate(eCtx, "dsb")
	assert.ErrorIs(t, err, errBadHash)
	err = echoFr.Close()
	require.NoError(t, err)
}

func TestAuthValidate(t *testing.T) {
	want := "abc"
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte(want)))
	eCtx := echoFr.AcquireContext()
	request.Header.Set(common.HashKeyHeaderName, "2f02e24ae2e1fe880399f27600afa88364e6062bf9bbe114b32fa8f23d03608a")
	eCtx.SetRequest(request)
	err := authValidate(eCtx, want)
	assert.NoError(t, err)
	err = echoFr.Close()
	require.NoError(t, err)
}

func TestAuthValidateBadHash(t *testing.T) {
	want := "abc"
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte(want)))
	eCtx := echoFr.AcquireContext()
	request.Header.Set(common.HashKeyHeaderName, "bad hash")
	eCtx.SetRequest(request)
	err := authValidate(eCtx, want)
	assert.ErrorIs(t, err, errBadHash)
	err = echoFr.Close()
	require.NoError(t, err)
}
