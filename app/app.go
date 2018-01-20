package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"io/ioutil"
	"net/http"
)

func getAppList(c echo.Context) error {

	result := []string{}

	files, err := ioutil.ReadDir(config.APP_DIR)

	if err != nil {
		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_SUCCESS,
			Result: result,
		})
	}

	for _, file := range files {
		if file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: result,
	})

}

type getAppStatusRequest struct {
	App string
}

func getAppStatus(c echo.Context) error {
	req := new(getAppStatusRequest)

	if err := c.Bind(&req); err != nil {

		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.REASON_PARAMETER_MISSING,
		})

	} else {

		return c.JSON(http.StatusOK, api.Status{
			Status: api.STATUS_SUCCESS,
			Result: module.GetAppStatus(req.App),
		})

	}
}

func RegisterHandlers(e *echo.Group) {

	e.Use(admin.AuthHandler)
	e.GET("/getAppList", getAppList)
	e.GET("/getAppStatus", getAppStatus)

}
