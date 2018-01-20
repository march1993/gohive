package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/march1993/gohive/app"
	"github.com/march1993/gohive/config"
)

func main() {
	fmt.Println("get_conf:" + config.Get("test", "test11"))

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	app.RegisterHandlers(e.Group("/app"))

	e.Logger.Fatal(e.Start(config.LISTEN_HOST_PORT))
}
