package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module/linux"
	"net/http"
)

type getAppListRequest struct {
	// empty
}

func getAppList(c echo.Context, request interface{}) error {
	// req := *request.(*getAppListRequest)

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: linux.GetAppList(),
		Addition: map[string]string{
			"PREFIX": config.APP_PREFIX,
		},
	})

}
