package main

import (
	"flag"

	"github.com/520lly/iamhere/app/controllers/api"
	"github.com/520lly/iamhere/app/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	//"github.com/labstack/gommon/log"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	//-----
	// API
	//-----
	// Echo instance
	api := echo.New()
	api.Debug = true
	api.Logger.SetPrefix("iamhere-api")

	api.Logger.Debug("Dialing mongo ", *mongo)
	api.Logger.Debug("Addr ", *addr)
	db.Init(*mongo, "", api.Logger)

	// Middleware
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())
	api.Use(middleware.BodyLimit("2K"))
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	//redirect HTTP to HTTPS
	//api.Pre(middleware.HTTPSRedirect())

	// Routes
	controllers.HandleMessages(api)
	controllers.HandleAreas(api)
	controllers.HandleAccounts(api)
	controllers.HandleTrail(api)

	//-----
	// WEB
	//-----
	//TBD

	// Start server
	api.Logger.Fatal(api.Start(*addr))
}
