package main

import (
	//"flag"
	//"io/ioutil"
	//"crypto/tls"
	"os"
	"time"

	"github.com/520lly/iamhere/app/controllers/api"
	"github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	//"golang.org/x/crypto/acme/autocert"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	//-----
	// API
	//-----
	// Echo instance
	api := echo.New()
	// Middleware
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	//load customized config
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
	controllers.HandleMessages(api)
	controllers.HandleAreas(api)
	controllers.HandleAccounts(api)
	controllers.HandleTrail(api)

	//-----
	// WEB
	//-----
	//TBD

	e := echo.New()
	//e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("www.historystest.com")
	// Cache certificates
	//e.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	controllers.HandleMessages(e)
	controllers.HandleAreas(e)
	controllers.HandleAccounts(e)
	controllers.HandleTrail(e)
	//go e.StartAutoTLS(":443")
	go e.StartTLS(":443", "/etc/ssl/214987401110045.pem", "/etc/ssl/214987401110045.key")

	// Start server
	api.Logger.Fatal(api.Start(Addr))
}
