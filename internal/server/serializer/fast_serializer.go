package serializer

import (
	"fmt"
	"net/http"
	"strings"

	goccyj "github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// FastJSONSerializer implements JSON encoding using encoding/json.
type FastJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d FastJSONSerializer) Serialize(c echo.Context, data interface{}, indent string) error {
	enc := goccyj.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}

	return enc.Encode(data) //nolint:wrapcheck
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d FastJSONSerializer) Deserialize(c echo.Context, data interface{}) error {
	err := goccyj.NewDecoder(c.Request().Body).Decode(data)
	if ute, ok := err.(*goccyj.UnmarshalTypeError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
			ute.Type, ute.Value, ute.Field, ute.Offset)

		return echo.NewHTTPError(http.StatusBadRequest, mess).SetInternal(err)
	} else if syne, ok := err.(*goccyj.SyntaxError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Syntax error: offset=%v, error=%v", syne.Offset, syne.Error())

		return echo.NewHTTPError(http.StatusBadRequest, mess).SetInternal(err)
	}

	return err //nolint:wrapcheck
}

func DeserializeString(body string, data interface{}) error {
	reader := strings.NewReader(body)
	err := goccyj.NewDecoder(reader).Decode(data)
	if ute, ok := err.(*goccyj.UnmarshalTypeError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v",
			ute.Type, ute.Value, ute.Field, ute.Offset)

		return fmt.Errorf("%s", mess)
	} else if syne, ok := err.(*goccyj.SyntaxError); ok { //nolint:errorlint
		mess := fmt.Sprintf("Syntax error: offset=%v, error=%v", syne.Offset, syne.Error())

		return fmt.Errorf("%s", mess)
	}

	return err
}
