package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"io/ioutil"
	"net/http"
)

type getAppListRequest struct {
	// empty
}

func getAppList(c echo.Context, request interface{}) error {

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
