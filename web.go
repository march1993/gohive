package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/march1993/gohive/app"
	"github.com/march1993/gohive/config"
)

func Web() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		/* TODO: CSRF */
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Static("/static", "./static")
	e.File("/favicon.ico", "static/favicon.ico")

	app.RegisterHandlers(e.Group("/app"))

	e.Logger.Fatal(e.Start(config.LISTEN_HOST_PORT))
}
