package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/golang"
	"net/http"
)

type getLastBuildRequest struct {
	App string
}

func getLastBuild(c echo.Context, request interface{}) error {
	req := *request.(*getLastBuildRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	return c.JSON(http.StatusOK, golang.GetLastBuild(req.App))

}
