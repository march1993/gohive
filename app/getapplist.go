package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module/linux"
	"io/ioutil"
	"net/http"
	. "strings"
)

type getAppListRequest struct {
	// empty
}

func getAppList(c echo.Context, request interface{}) error {
	// req := *request.(*getAppListRequest)

	result := []string{}

	files, err := ioutil.ReadDir(config.APP_DIR)

	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		if file.IsDir() {
			name := file.Name()

			if !HasSuffix(name, linux.Suffix) {
				result = append(result, name)
			}
		}
	}

	return c.JSON(http.StatusOK, api.Status{
		Status: api.STATUS_SUCCESS,
		Result: result,
	})

}
