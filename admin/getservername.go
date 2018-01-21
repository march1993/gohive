package admin

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"net/http"
)

type getServerNameRequest struct {
}

func getServerName(c echo.Context, request interface{}) error {

	serverName := config.Get("server_name", "_")
	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: serverName,
	})

}
