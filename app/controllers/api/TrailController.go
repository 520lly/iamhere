package controllers

import (
	. "github.com/520lly/iamhere/app/iamhere"
	"github.com/labstack/echo"
)

func HandleTrail(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Trails.Group
	g := e.Group(urlGroup)
	g.GET("/", GetMessages)
	g.GET("", GetMessages)
}
