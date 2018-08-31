package iamhere

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Log(msg string, i ...interface{}) {
	echo.New().Use(middleware.Logger())
}
