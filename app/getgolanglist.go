package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module/golang"
	"net/http"
)

type getGolangListRequest struct {
}

func getGolangList(c echo.Context, request interface{}) error {
	// req := *request.(*getGolangListRequest)

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: golang.GetGolangList(),
	})

}
