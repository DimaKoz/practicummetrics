package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DimaKoz/practicummetrics/internal/common"
	"github.com/DimaKoz/practicummetrics/internal/common/config"
	"github.com/DimaKoz/practicummetrics/internal/common/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRsaAesDecoderNoRsaEncodedKeyHeader(t *testing.T) {
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte("")))
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)
	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := rsaAesD(ctx)
	assert.Error(t, err)
}

func TestRsaAesDecoderNoBody(t *testing.T) {
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(nil))

	request.Header.Add(common.RsaEncodedKeyHeaderName, "qwerty")
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)
	ctx.Request().Body = nil
	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := rsaAesD(ctx)
	assert.Error(t, err)
}

func TestRsaAesDecoderNoBody1(t *testing.T) {
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(nil))

	request.Header.Add(common.RsaEncodedKeyHeaderName, "qwerty")
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)

	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := rsaAesD(ctx)
	assert.Error(t, err)
}

func TestRsaAesDecoderBadEncoded(t *testing.T) {
	//nolint:exhaustruct
	cfg := config.ServerConfig{Config: config.Config{CryptoKey: "keys/keyfile.pem"}}
	repository.LoadPrivateKey(cfg)
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader([]byte("hop")))

	request.Header.Add(common.RsaEncodedKeyHeaderName, "qwerty")
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)

	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := rsaAesD(ctx)
	assert.Error(t, err)
}

var (
	rsaEncoded = "MdJuIQ77YjZcLCA9HlFkY8YrUfpoXUK64sFvChXKt9HK09bG171trkvGMSkq1loaSlXfk8C0o326T0IzwMvvmfVNdayHzlA5" +
		"WSPi2ECTbA9ThC/lVmsbGcId+nLdDXPqMbHZweDUZrCqW29gFTBlbOAa7rWZa80s0ChkyarF9dpZnlk3GhdE+RnaSztY0CQNPnqbKxWVP" +
		"MYkCqjPEso/UpBaHlm1JyRacGPIrwCk8ihEzkONRyp" +
		"8zoeKjJmCJsxZH1h5rK+9nAdJkDoy0ZESi0fyEKd0K+Pm5wR06VxlLBznqlZIuHeeOpmqMvoWyaluvCnWs1nfow3fZ7/JevpC8Q"
	encMess = []byte("QDurHTxozRN8weYd8qdYRq4j8Ly4PZR7NvX7ZBvUNuo5ocsBCjyIYvxLxaQ3qwe" +
		"/3m9IaQAMfqu" +
		"RuM3GO85t4hu2CVgixxQghXT3TlDYhDNb52r0wtCkoIYVpVby0WuboAYmfI" +
		"/5yxsCZLH7l8IYJBHMq7c/" +
		"3kxz73bvYPvDbf5Rm5j020jDYZ3xTzPBMgN4CzbJRpSUAALa" +
		"/ZuWdvo9mxb7k+" +
		"obdR2AoQ1yYgIIIvfqxV9BdJ5sdo4MCpq4" +
		"/VzFjmHruMDov27U0FOVyluaIRJGZdAqVSic622wA5b8A04Jw1IXYkJNWLN" +
		"/vA6gSmaaM2vz/b6Ugsi3uTiqlq+WI6asxeQEI0N0nTIgJpxXL4ua93n0vs4geMdYsgLEay3p9tVQeklIDjztFnLkJmvLb" +
		"/r9axW+G7oNb8V2voboqo7vA82fDwe+QIEJPgXqcpf1vOaJ03durn2P" +
		"/1BJi4JzZP6AiFVWcTm0lXfbVPruo12l5b9AFDwQi7zU4svlt+de8sg7g8ht8HesuuKkHte" +
		"/MNQG2cK86bRPq7kg+FJAyeJryqR8ru3EAp2C+NM49l/nPUEypN6s0K+U+QyYYYm8iE" +
		"/hyLObG2OnD97KkOImIr6izikTI" +
		"/eIg6diyK91o+V+hB5QgnTKjVEXOf5Q33ok7v515KqyEF3MhGlPXpkDd6pyK4HLU21+16GbC7EERnCeLCRnF" +
		"/O7axgCHsK3Qzl2rtw9ggNowjugFCWYQV1Qu3a/kurgkiortaDe0RV94wY+Ll8mDeYhkslq0OU6i6bEofTbj6D7QZNKJwb7UVR2jEoPzCf+" +
		"/ZVflaD5hUBqDTgffM" +
		"/lduWxj6gCRUhipAWriVItiEEEifSbg6RrtpQyI54KAtZ4nRowq8U6Avszrex1gwdA9umpqKLc5TfE35u84YMjwXOQxggCm6s4Mb3NWHZ5OiR" +
		"/EWuT4tdE4CKV8RtGimg1fdR6vvJsSCm5mot5pWPUpNEzFnyn9xCsOoOCNmC5C+QLWRxPX4S0p3FDi9Xt1q9abczunF+2GOkxllqgsfSbzV" +
		"/ZLTFp/iUjdFIfRaQHxMX4C5B+A5Q6ZvwH6N2lS8++QTqpAMDRp52dXpApjEtvPc" +
		"/pCl8FTcj5N8O6K6tGZUkjT3BVS5cEYlz3lo82R2sfRJZA3andnMMarkN2MKe4XvMhamXOsGecU1o4NQ+Q5p0w" +
		"/NMaFRQ8c3HjGEDiAwh7WHfvb4XwvnNQ2FHkqNGzQXVkXy8mO2rM3gnvJ9W4rLCxte4JlG9JoJv5mY" +
		"rBee45L6JlNO5PTy3knnycGmGPQ3v5ptNBcpiljMiU9NT1LW0HIRwrsv2FxXQndv" +
		"/k8moSKgiv5WgralvCdAzKchHkGLbrVLGHkKG3i9Wk4GylmcidYU8BZqy8Kt3" +
		"/8HvYeodg9bU8+AJJnJUdcpy4EeT0IwT9K+kKeudXcoS8G40iKypzDLcfu6AK0JpPhMHWTwuwu8SDaE71JviI3n5NSXZkQQXt8pUthd1QiW" +
		"uYSfSPcIohElH5R2YAz6kV2ozX2dWq48ktDhpe6PP9iboPm6pRwRtzwhOyQyLRfYxFOYjY5gCxiuSwbcIv8J3ezZa7awPWckJs10nlz4r" +
		"/YMe4n2qefoKyTo6GyXRN/6tKjultwiXp5ez1+KkbJl2o6c/K8AWTbHsr49GPqK2md+nIqeP07HKu6sDhuDKYzcrxL6Fhf+DLN3rupvrcoLg" +
		"QswZhlf5uTbKMhhx65sFhYt0Y5O8dmj72Goex+ORnuPg1ePWHgvmJ+UkiYyTWlxq3/+y+uBnZbrVx5DQhwHvyPEYh+" +
		"FezUBVcGY0vyjjMYaq7H+dXWbhgM7U6fQgeP0zQKZrDuwU5IC2xm0anYfMUxXoF9B2WxpPa0Ze4A1NJt+8KX9YOIOpOS9sYAPDVocXm9su" +
		"N5OwVC3vixkkwXjs7w73LmCYuON/wZUZvGK5ovq9B+KXabAv69j+QkQtKILwEsLtHwLC9p1gMR8r8GFg4Fwzlo9KD9M7KqWloTUhABCDZz2" +
		"QFO7ZbDvHV7KfsdNamWq2dpcpqaIc1hSG9cfYsBo9uUzShHw")
)

func TestRsaAesDecoderEncoded(t *testing.T) {
	//nolint:exhaustruct
	cfg := config.ServerConfig{Config: config.Config{CryptoKey: "keys/keyfile.pem"}}
	repository.LoadPrivateKey(cfg)
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(encMess))

	request.Header.Add(common.RsaEncodedKeyHeaderName, rsaEncoded)
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)

	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		//nolint:wrapcheck
		return c.NoContent(http.StatusOK)
	})
	err := rsaAesD(ctx)
	assert.NoError(t, err)
}

func TestRsaAesDecoderErrNext(t *testing.T) {
	//nolint:exhaustruct
	cfg := config.ServerConfig{Config: config.Config{CryptoKey: "keys/keyfile.pem"}}
	repository.LoadPrivateKey(cfg)
	echoFr := echo.New()
	request := httptest.NewRequest(http.MethodGet, "https://example.com", bytes.NewReader(encMess))

	request.Header.Add(common.RsaEncodedKeyHeaderName, rsaEncoded)
	recorder := httptest.NewRecorder()
	ctx := echoFr.NewContext(request, recorder)

	rsaAesD := RsaAesDecoder()(func(c echo.Context) error {
		return http.ErrBodyNotAllowed
	})
	err := rsaAesD(ctx)
	assert.ErrorIs(t, err, http.ErrBodyNotAllowed)
}
