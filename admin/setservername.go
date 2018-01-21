package admin

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module/nginx"
	"net/http"
)

type setServerNameRequest struct {
	ServerName string
}

func setServerName(c echo.Context, request interface{}) error {
	req := *request.(*setServerNameRequest)

	if checkServerName(req.ServerName) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.SERVER_NAME_INVALID,
		})
	}

	config.Set("server_name", req.ServerName)

	nginx.RegisterNginx()

	return c.JSON(http.StatusOK, api.Status{Status: api.STATUS_SUCCESS})

}
