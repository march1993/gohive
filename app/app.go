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

	e.POST("/statusApp", api.EnsureRequest(statusApp, &statusAppRequest{}))
	e.POST("/createApp", api.EnsureRequest(createApp, &createAppRequest{}))
	e.POST("/repairApp", api.EnsureRequest(repairApp, &repairAppRequest{}))
	e.POST("/removeApp", api.EnsureRequest(removeApp, &removeAppRequest{}))
	e.POST("/renameApp", api.EnsureRequest(renameApp, &renameAppRequest{}))
	e.POST("/listRemovedApp", api.EnsureRequest(listRemovedApp, &listRemovedAppRequest{}))

	e.POST("/getGitUrl", api.EnsureRequest(getGitUrl, &getGitUrlRequest{}))
	e.POST("/setGitKeys", api.EnsureRequest(setGitKeys, &setGitKeysRequest{}))
	e.POST("/getGitKeys", api.EnsureRequest(getGitKeys, &getGitKeysRequest{}))

}
