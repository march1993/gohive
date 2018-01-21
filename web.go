package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/app"
	"github.com/march1993/gohive/config"
	"net/http"
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
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	app.RegisterHandlers(e.Group("/app"))
	admin.RegisterHandlers(e.Group("/admin"))

	e.Logger.Fatal(e.Start(config.LISTEN_HOST_PORT))
}
