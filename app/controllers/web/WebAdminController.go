package controllers_web

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strings"

	. "github.com/520lly/iamhere/app/iamhere"
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

var urlGroup string

func HandleWebAdmin(e *echo.Echo) {
	urlGroup = Config.WebConfig.Prefix + Config.WebConfig.Version + Config.WebConfig.Admin.Group
	g := e.Group(urlGroup)
	g.Static("/css/", "web/template")
	g.Static("/js/", "web/template")
	g.Static("/img/login/", "web/template")
	//g.Static("/html/", "web/template")
	//g.File("", "web/template/html/index.html")
	g.GET("", AdminHandler)
	g.GET("/login/", urlHandler)
	//g.Use(middleware.Static("web/template/html"))
	//g.Use(middleware.Static("web/public/"))
	//g.Static("css", http.FileServer(http.Dir("web/template")))
	//设置静态资源
	//设置路由
	//http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/ajax/", ajaxHandler)
	http.HandleFunc("/getUsers/", ajaxHandler)
	http.HandleFunc("/", NotFoundHandler)
	//g.GET("", ValidateAccount)
	//g.POST("/login", ValidateAccount)
	//g.GET("/login/", ValidateAccount)
}

func AdminHandler(e echo.Context) error {
	// 获取cookie
	cookie, err := e.Request().Cookie("username")
	if err != nil || cookie.Value == "" {
		log.Println("重定向到登录界面")
		e.Redirect(http.StatusFound, urlGroup+"/login")
		return err
	}
	return nil

	pathInfo := strings.Trim(e.Path(), "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "View"
	} else {
		action = strings.Title(parts[0] + "View")
	}

	admin := &URLController{}
	controller := reflect.ValueOf(admin)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "View")
	}
	requestValue := reflect.ValueOf(e.Request())
	responseValue := reflect.ValueOf(e.Response())
	userValue := reflect.ValueOf(cookie.Value)
	method.Call([]reflect.Value{responseValue, requestValue, userValue})
	return nil
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	t, err := template.ParseFiles("template/html/404.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
}

func urlHandler(e echo.Context) error {
	pathInfo := strings.Trim(e.Path(), "/")
	parts := strings.Split(pathInfo, "/")
	log.Println(pathInfo)
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "View"
	} else {
		action = strings.Title(parts[0]) + "View"
	}

	login := &URLController{}
	controller := reflect.ValueOf(login)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "View")
	}
	log.Println(action)
	requestValue := reflect.ValueOf(e.Request())
	responseValue := reflect.ValueOf(e.Response())
	method.Call([]reflect.Value{responseValue, requestValue})
	return nil
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	pathInfo := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(pathInfo, "/")
	var action = ""
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	} else {
		action = strings.Title(parts[0]) + "Action"
	}
	ajax := &ajaxController{}
	controller := reflect.ValueOf(ajax)
	method := controller.MethodByName(action)
	if !method.IsValid() {
		method = controller.MethodByName(strings.Title("index") + "Action")
	}
	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}
