package main

import (
	//"net/http"
	"os"
	"time"

	. "github.com/520lly/iamhere/app/controllers/api"
	. "github.com/520lly/iamhere/app/controllers/web"
	"github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	// Hosts
	hosts := map[string]*Host{}
	//-----
	// API
	//-----
	// Echo instance
	api := echo.New()
	// Middleware
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	InitConfig(api.Logger)
	if Config.AppConfig.EnableDebug {
		api.Logger.SetLevel(log.DEBUG)
	}
	if Config.AppConfig.EnableDumpLog {
		// set location of log file
		now := time.Now()
		dt := now.Format("2006-01-02-15:04:05")
		var logpath = Config.AppConfig.LogPath + dt + ".log"
		api.Logger.Debug("Dump Log file: ", logpath)
		f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			api.Logger.Fatal("error opening file: ", err)
		}
		defer f.Close()
		api.Logger.SetOutput(f)
	}
	api.Logger.SetPrefix(Config.AppConfig.LoggerPrefix)
	api.Logger.Debug("Dialing mongo ", Config.Database.Host)
	Addr := string(Config.AppConfig.Host) + ":" + string(Config.AppConfig.Port)
	api.Logger.Debug("Addr ", Addr)
	db.Init(Config.Database.Host, Config.Database.Name, api.Logger)

	api.Use(middleware.BodyLimit(Config.ApiConfig.BodySizeLimit))
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	//redirect HTTP to HTTPS
	//api.Pre(middleware.HTTPSRedirect())

	// Routes
	HandleMessages(api)
	HandleAreas(api)
	HandleAccounts(api)
	HandleTrail(api)

	hosts[Addr] = &Host{api}

	//-----
	// WEB
	//-----
	site := echo.New()
	// Middleware
	site.Use(middleware.Logger())
	site.Use(middleware.Recover())
	hosts["site."+Addr] = &Host{site}

	HandleWebAdmin(site)

	// Server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]

		if host == nil {
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})
	// Start server
	//api.Logger.Fatal(api.Start(Addr))
	e.Logger.Fatal(e.Start(Addr))
}
