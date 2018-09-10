package controllers_web

import (
	"github.com/labstack/echo"
	"html/template"

	. "github.com/520lly/iamhere/app/modules"
)

type URLController struct {
}

func (this *URLController) IndexView(c echo.Context, user string) {
	t, err := template.ParseFiles("template/html/index.html")
	if err != nil {
		c.Logger().Error(err)
	}
	t.Execute(c.Response().Writer, &AdminUser{user})
}

func (this *URLController) LoginView(c echo.Context) {
	t, err := template.ParseFiles("template/html/login.html")
	if err != nil {
		c.Logger().Error(err)
	}
	t.Execute(c.Response().Writer, nil)
}
