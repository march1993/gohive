package admin

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
)

func RegisterHandlers(e *echo.Group) {

	e.Use(AuthHandler)

	e.POST("/setServerName", api.EnsureRequest(setServerName, &setServerNameRequest{}))
	e.POST("/getServerName", api.EnsureRequest(getServerName, &getServerNameRequest{}))

}
