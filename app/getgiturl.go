package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/git"
	"net/http"
)

type getGitUrlRequest struct {
	AppRequest
}

func getGitUrl(c echo.Context, request interface{}) error {
	req := *request.(*getGitUrlRequest)

	if checkName(req.App) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	gitUrl := git.GetGitUrl(req.App)
	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: gitUrl,
		Info:   "git push " + gitUrl + " <branch>",
	})

}
