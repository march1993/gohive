package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/golang"
	"net/http"
)

type setGolangVersionRequest struct {
	App     string
	Version string
}

func setGolangVersion(c echo.Context, request interface{}) error {
	req := *request.(*setGolangVersionRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	if !golang.CheckGolangVersion(req.Version) {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.GOLANG_VERSION_INVALID,
		})
	}

	golang.SetGolangVersion(req.App, req.Version)

	return repairApp(c, &repairAppRequest{App: req.App})

}
