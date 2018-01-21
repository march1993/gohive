package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/golang"
	"net/http"
)

type setGolangInstallationRequest struct {
	Version string
}

func setGolangInstallation(c echo.Context, request interface{}) error {
	req := *request.(*setGolangInstallationRequest)

	if !golang.CheckGolangVersion(req.Version) {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.GOLANG_VERSION_INVALID,
		})
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: golang.SetGolangInstallation(req.Version),
	})

}
