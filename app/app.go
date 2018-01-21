package app

import (
	"github.com/labstack/echo"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"os"
)

func init() {
	if err := os.MkdirAll(config.APP_DIR, 0755); err != nil {
		panic(err.Error())
	}
}

func RegisterHandlers(e *echo.Group) {

	e.Use(admin.AuthHandler)
	e.POST("/getAppList", api.EnsureRequest(getAppList, &getAppListRequest{}))

	e.POST("/getAppStatus", api.EnsureRequest(getAppStatus, &getAppStatusRequest{}))
	e.POST("/createApp", api.EnsureRequest(createApp, &createAppRequest{}))
	e.POST("/repairApp", api.EnsureRequest(repairApp, &repairAppRequest{}))

}
