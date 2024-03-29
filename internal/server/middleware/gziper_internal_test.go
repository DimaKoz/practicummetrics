package middleware

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestGzipSkipperHasGzip(t *testing.T) {
	echoFramework := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	responseRecorder := httptest.NewRecorder()

	request.Header.Set(echo.HeaderAcceptEncoding, "gzip")
	ctx := echoFramework.NewContext(request, responseRecorder)
	got := gzipSkipper(ctx)
	assert.False(t, got)
	contentEnc := ctx.Response().Header().Get(echo.HeaderContentEncoding)
	assert.Contains(t, contentEnc, "gzip")
}

func TestGzipSkipperNoGzip(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(request, w)
	got := gzipSkipper(ctx)
	assert.True(t, got)
	contentEnc := ctx.Response().Header().Get(echo.HeaderContentEncoding)
	assert.Empty(t, contentEnc)
}

func TestNewGzipConfig(t *testing.T) {
	type args struct {
		f func(c echo.Context) bool
	}
	tests := []struct {
		name string
		args args
		want middleware.GzipConfig
	}{
		{
			name: "nil Skipper",
			args: args{f: nil},
			want: middleware.GzipConfig{
				Skipper:   nil,
				Level:     gzip.BestCompression,
				MinLength: 0,
			},
		},
	}
	for _, test := range tests {
		testUnit := test
		t.Run(testUnit.name, func(t *testing.T) {
			assert.Equalf(t, testUnit.want, newGzipConfig(testUnit.args.f), "newGzipConfig(%v)", testUnit.args.f)
		})
	}
}

//nolint:exhaustruct
func TestGetGzipMiddlewareConfig(t *testing.T) {
	tests := []struct {
		name string
		want echo.MiddlewareFunc
	}{
		{
			want: middleware.GzipWithConfig(middleware.GzipConfig{
				Skipper: nil,
				Level:   gzip.BestCompression,
			}),
		},
	}

	for _, testItem := range tests {
		test := testItem
		t.Run(test.name, func(t *testing.T) {
			orig := gzipSkipper
			gzipSkipper = nil
			t.Cleanup(func() {
				gzipSkipper = orig
			})
			assert.NotNil(t, test.want, GetGzipMiddlewareConfig(), "GetGzipMiddlewareConfig()")
		})
	}
}
