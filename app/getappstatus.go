package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"net/http"
)

type getAppStatusRequest struct {
	App string
}

func getAppStatus(c echo.Context, request interface{}) error {
	req := *request.(*getAppStatusRequest)

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: module.GetAppStatus(req.App),
	})

}
