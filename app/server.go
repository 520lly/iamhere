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
	"golang.org/x/crypto/acme/autocert"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func main() {
	//-----
	// HTTP API
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
		var logpath = Config.AppConfig.LogPath + "HTTP-" + dt + ".log"
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

	api.Use(middleware.BodyLimit(Config.ApiConfig.BodySizeLimit))
	api.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	//redirect HTTP to HTTPS
	//api.Pre(middleware.HTTPSRedirect())

	api.Logger.Debug("JWT secret ", GetJWTSecretCode())
	// Routes
	controllers.HandleMessages(api)
	controllers.HandleAreas(api)
	controllers.HandleAccounts(api)
	controllers.HandleLogin(api)
	controllers.HandleTrail(api)

	//-----
	// API HTTPS
	//-----
	apiTLS := echo.New()
	apiTLS.Use(middleware.Logger())
	apiTLS.Use(middleware.Recover())

	//load customized config
	InitConfig(apiTLS.Logger)
	if Config.AppConfig.EnableDebug {
		apiTLS.Logger.SetLevel(log.DEBUG)
	}
	apiTLS.Logger.SetPrefix(Config.AppConfig.LoggerPrefix)
	if Config.AppConfig.EnableDumpLog {
		// set location of log file
		now := time.Now()
		dt := now.Format("2006-01-02-15:04:05")
		var logpath = Config.AppConfig.LogPath + "HTTPS-" + dt + ".log"
		apiTLS.Logger.Debug("Dump Log file: ", logpath)
		f, err := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			apiTLS.Logger.Fatal("error opening file: ", err)
		}
		defer f.Close()
		apiTLS.Logger.SetOutput(f)
	}
	apiTLS.Logger.SetPrefix(Config.AppConfig.LoggerPrefix)
	db.Init(Config.Database.Host, Config.Database.Name, apiTLS.Logger)

	controllers.HandleMessages(apiTLS)
	controllers.HandleAreas(apiTLS)
	controllers.HandleAccounts(apiTLS)
	controllers.HandleLogin(apiTLS)
	controllers.HandleTrail(apiTLS)

	apiTLS.Logger.Debug("enableSSL ", Config.AppConfig.EnableSSL)
	if Config.AppConfig.EnableSSL {
		apiTLS.AutoTLSManager.HostPolicy = autocert.HostWhitelist("www.historystest.com")
		//Cache certificates
		apiTLS.AutoTLSManager.Cache = autocert.DirCache("/etc/letsencrypt/live/www.historystest.com")
		go apiTLS.StartAutoTLS(":443")
	}

	// Start server
	api.Logger.Fatal(api.Start(Addr))
}
