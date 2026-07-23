package artifacts

import (
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func extractObjectKeyFromPath(c echo.Context) (string, error) {
	handlerPath := strings.TrimSuffix(c.Path(), "*")
	objKey := strings.Replace(c.Request().URL.Path, handlerPath, "", 1)
	unescapedObjKey, err := url.QueryUnescape(objKey)
	if err != nil {
		return "", err
	}
	return unescapedObjKey, nil
}
