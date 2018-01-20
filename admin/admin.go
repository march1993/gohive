package admin

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"net/http"
)

func TestToken(token string) bool {
	return token != "" && config.Get("token", "") == token
}
func SetToken(token string) {
	config.Set("token", token)
}

func AuthHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		credential := new(api.Credential)

		if err := c.Bind(&credential); err != nil {
			// try get token from headers
			header := c.Request().Header
			credential.Token = header.Get("Token")
		}

		if TestToken(credential.Token) {
			return next(c)
		} else {
			return c.JSON(http.StatusOK, api.Status{Status: api.STATUS_FAILURE, Reason: api.AUTH_FAILURE})
		}

	}
}
