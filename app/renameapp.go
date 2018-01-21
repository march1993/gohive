package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"net/http"
)

type renameAppRequest struct {
	OldName string
	NewName string
}

func renameApp(c echo.Context, request interface{}) error {
	req := *request.(*renameAppRequest)

	if checkName(req.NewName) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	if checkName(req.OldName) == false {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.APP_NAME_INVALID,
		})
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: module.RenameApp(req.OldName, req.NewName),
	})

}
