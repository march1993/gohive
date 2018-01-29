package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/git"
	"net/http"
)

type setGitKeysRequest struct {
	App  string
	Keys []string
}

func setGitKeys(c echo.Context, request interface{}) error {
	req := *request.(*setGitKeysRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	return c.JSON(http.StatusOK, git.SetGitKeys(req.App, req.Keys))

}
