package iamhere

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"net/http"
)

func DecodeBody(c echo.Context, v interface{}) error {
	defer c.Request().Body.Close()
	return json.NewDecoder(c.Request().Body).Decode(v)
}

func EncodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func RespondJ(c echo.Context, status int, data interface{}) error {
	if data != nil {
		return c.JSON(status, data)
	}
	return NewError("data is nil")
}

func RespondS(c echo.Context, status int, data string) error {
	return c.String(status, data)
}

func RespondErr(c echo.Context, status int, args ...interface{},
) {
	RespondJ(c, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

//func RespondHTTPErr(w http.ResponseWriter, r *http.Request,
//    status int,
//) {
//    RespondErr(w, r, status, http.StatusText(status))
//}
