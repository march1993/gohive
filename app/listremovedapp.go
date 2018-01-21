package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"net/http"
)

type listRemovedAppRequest struct {
	AppRequest
}

func listRemovedApp(c echo.Context, request interface{}) error {
	// req := *request.(*listRemovedAppRequest)

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: module.ListRemovedApp(),
	})

}
