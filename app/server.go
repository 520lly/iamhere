package main

import (
	"flag"

	"github.com/520lly/iamhere/app/controllers"
	"github.com/520lly/iamhere/app/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	// Echo instance
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	e.Logger.Debug("Dialing mongo ", *mongo)
	e.Logger.Debug("Addr ", *addr)
	db.Init(*mongo, "")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Routes
	controllers.HandleMessages(e)
	controllers.HandleAreas(e)
	// Start server
	e.Logger.Fatal(e.Start(*addr))
}
