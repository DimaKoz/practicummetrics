package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_gzipSkipperHasGzip(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.Header.Set(echo.HeaderAcceptEncoding, "gzip")
	w := httptest.NewRecorder()
	ctx := e.NewContext(request, w)
	got := gzipSkipper(ctx)
	assert.False(t, got)
	contentEnc := ctx.Response().Header().Get(echo.HeaderContentEncoding)
	assert.Contains(t, contentEnc, "gzip")
}

func Test_gzipSkipperNoGzip(t *testing.T) {
	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()
	ctx := e.NewContext(request, w)
	got := gzipSkipper(ctx)
	assert.True(t, got)
	contentEnc := ctx.Response().Header().Get(echo.HeaderContentEncoding)
	assert.Empty(t, contentEnc)
}

func Test_newGzipConfig(t *testing.T) {
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
				Skipper: nil,
				Level:   5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, newGzipConfig(tt.args.f), "newGzipConfig(%v)", tt.args.f)
		})
	}
}

func TestGetGzipMiddlewareConfig(t *testing.T) {
	tests := []struct {
		name string
		want echo.MiddlewareFunc
	}{
		{
			want: middleware.GzipWithConfig(middleware.GzipConfig{
				Skipper: nil,
				Level:   5,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := gzipSkipper
			gzipSkipper = nil
			t.Cleanup(func() {
				gzipSkipper = orig
			})
			assert.NotNil(t, tt.want, GetGzipMiddlewareConfig(), "GetGzipMiddlewareConfig()")
		})
	}
}
