package controllers_web

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"

	. "github.com/520lly/iamhere/app/iamhere"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var urlGroup string

func HandleWebAdmin(e *echo.Echo) {
	urlGroup = Config.WebConfig.Prefix + Config.WebConfig.Version + Config.WebConfig.Admin.Group
	g := e.Group(urlGroup)
	e.Use(middleware.Static("template"))
	g.GET("/admin", AdminHandler)
	g.GET("/login", urlHandler)
	g.POST("/ajax/login", ajaxHandler)
	g.GET("/getUsers", ajaxHandler)
	g.GET("/", NotFoundHandler)
}

func AdminHandler(c echo.Context) error {
	// 获取cookie
	cookie, err := c.Cookie("username")
	if err != nil || cookie.Value == "" {
		c.Logger().Debug("重定向到登录界面")
		c.Redirect(http.StatusFound, urlGroup+"/login")
	}

	pathInfo := strings.Trim(c.Path(), "/")
	c.Logger().Debug("pathInfo:", pathInfo)
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "View"
	} else {
		action = strings.Title(parts[0] + "View")
	}

	c.Logger().Debug("action:", action)
	admin := &URLController{}
	controller := reflect.ValueOf(admin)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "View")
	}
	username := "root"
	c.Logger().Debug("method:", &method)
	echoValue := reflect.ValueOf(c)
	userValue := reflect.ValueOf(username)
	method.Call([]reflect.Value{echoValue, userValue})
	return nil
}

func NotFoundHandler(c echo.Context) error {
	if c.Path() == "/" {
		c.Redirect(http.StatusFound, urlGroup+"/login")
	}

	t, err := template.ParseFiles("web/template/html/404.html")
	if err != nil {
		log.Println(err)
		return err
	}
	t.Execute(c.Response().Writer, nil)
	return nil
}

func urlHandler(c echo.Context) error {
	pathInfo := strings.Trim(c.Path(), "/")
	c.Logger().Debug("pathInfo:", pathInfo)
	parts := strings.Split(pathInfo, "/")
	c.Logger().Debug("pathInfo:", pathInfo)
	c.Logger().Debug("parts:", parts)
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "View"
	} else {
		action = strings.Title(parts[0]) + "View"
	}

	c.Logger().Debug("action:", action)
	login := &URLController{}
	controller := reflect.ValueOf(login)
	method := controller.MethodByName(action)
	c.Logger().Debug("method:", method)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "View")
	}
	c.Logger().Debug("method:", &method)
	echoValue := reflect.ValueOf(c)
	//username := "root"
	//userValue := reflect.ValueOf(username)
	method.Call([]reflect.Value{echoValue})
	return nil
}

func ajaxHandler(c echo.Context) error {
	pathInfo := strings.Trim(c.Path(), "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	} else {
		action = strings.Title(parts[0]) + "Action"
	}
	c.Logger().Debug("action:", action)
	ajax := &AjaxController{}
	controller := reflect.ValueOf(ajax)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	c.Logger().Debug("method:", &method)
	echoValue := reflect.ValueOf(c)
	method.Call([]reflect.Value{echoValue})
	return nil
}
