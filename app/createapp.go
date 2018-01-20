package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"net/http"
)

type createAppRequest struct {
	App string
}

func createApp(c echo.Context, request interface{}) error {
	req := *request.(*createAppRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	for _, module := range module.ModuleList {
		if err := module.Create(req.App); err != nil {
			return c.JSON(http.StatusOK, api.Status{
				Status: api.STATUS_FAILURE,
				Reason: err.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
	})

}
