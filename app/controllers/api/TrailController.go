package controllers

import (
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/labstack/echo"
)

func HandleTrail(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Trails.Group
	g := e.Group(urlGroup)
	g.GET("", GetTrailMessages)
	g.GET("/:longitude&:latitude", GetMessages)
}

func GetTrailMessages(c echo.Context) error {
	var msg Message
	debugF := true
	if err := HandleGetMessages(c, &msg, debugF); err != nil {
		return err
	}
	return nil
}
