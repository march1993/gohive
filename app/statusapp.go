package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"net/http"
)

type statusAppRequest struct {
	App string
}

func statusApp(c echo.Context, request interface{}) error {
	req := *request.(*statusAppRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: module.StatusApp(req.App),
	})

}
