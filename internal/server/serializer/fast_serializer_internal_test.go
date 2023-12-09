package serializer

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	goccyj "github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
	assrt "github.com/stretchr/testify/assert"
)

const (
	userJSON       = `{"id":1,"name":"Jon Snow"}`
	invalidContent = "invalid content"
)

const userJSONPretty = `{
  "id": 1,
  "name": "Jon Snow"
}`

type (
	user struct {
		//nolint:tagliatelle
		ID int `form:"id" header:"id" json:"id" param:"id" query:"id" xml:"id"`
		//nolint:tagliatelle
		Name string `form:"name" header:"name" json:"name" param:"name" query:"name" xml:"name"`
	}
)

func TestFastJSONSerializerSerialize(t *testing.T) {
	assert := assrt.New(t)
	echoFr := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctxEcho := echoFr.NewContext(req, rec)

	// Echo
	assert.Equal(echoFr, ctxEcho.Echo())

	// Request
	assert.NotNil(ctxEcho.Request())

	// Response
	assert.NotNil(ctxEcho.Response())

	//--------
	// FastJSONSerializer JSON encoder
	//--------

	enc := new(FastJSONSerializer)

	err := enc.Serialize(ctxEcho, user{1, "Jon Snow"}, "")
	if assert.NoError(err) {
		assert.Equal(userJSON+"\n", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodPost, "/", nil)
	rec = httptest.NewRecorder()
	ctxEcho = echoFr.NewContext(req, rec)
	err = enc.Serialize(ctxEcho, user{1, "Jon Snow"}, "  ")
	if assert.NoError(err) {
		assert.Equal(userJSONPretty+"\n", rec.Body.String())
	}
}

func TestFastJSONSerializerDecode(t *testing.T) {
	echoFr := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()
	ctxEcho := echoFr.NewContext(req, rec)

	assert := assrt.New(t)

	// Echo
	assert.Equal(echoFr, ctxEcho.Echo())

	// Request
	assert.NotNil(ctxEcho.Request())

	// Response
	assert.NotNil(ctxEcho.Response())

	//--------
	// FastJSONSerializer JSON encoder
	//--------

	enc := new(FastJSONSerializer)

	//nolint:exhaustruct
	u := user{}
	err := enc.Deserialize(ctxEcho, &u)
	if assert.NoError(err) {
		assert.Equal(u, user{ID: 1, Name: "Jon Snow"})
	}

	//nolint:exhaustruct
	userUnmarshalSyntaxError := user{}
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(invalidContent))
	rec = httptest.NewRecorder()
	ctxEcho = echoFr.NewContext(req, rec)
	err = enc.Deserialize(ctxEcho, &userUnmarshalSyntaxError)
	//nolint:exhaustruct
	assert.IsType(&echo.HTTPError{}, err)
	//nolint:errorlint, forcetypeassert
	errInternal := (err.(*echo.HTTPError)).Unwrap()
	//nolint:exhaustruct
	assert.IsType(&goccyj.SyntaxError{}, errInternal)

	//nolint:exhaustruct
	userUnmarshalTypeError := struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{}

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec = httptest.NewRecorder()
	ctxEcho = echoFr.NewContext(req, rec)
	err = enc.Deserialize(ctxEcho, &userUnmarshalTypeError)
	//nolint:exhaustruct
	assert.IsType(&echo.HTTPError{}, err)
	//nolint:errorlint,forcetypeassert
	errInternal = (err.(*echo.HTTPError)).Unwrap()
	//nolint:exhaustruct
	assert.IsType(&goccyj.UnmarshalTypeError{}, errInternal)
}
