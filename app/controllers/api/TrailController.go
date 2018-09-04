package controllers

import (
	"github.com/labstack/echo"
)

func HandleTrail(e *echo.Echo) {
	g := e.Group("/trail")
	g.GET("/", GetMessages)
	g.GET("", GetMessages)
}
