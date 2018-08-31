package controllers

import (
	"net/http"
	//"strconv"
	//"time"

	//"github.com/520lly/iamhere/app/db"
	"github.com/labstack/echo"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

// Handlers
func CreateNewMessage(c echo.Context) error {

	return c.String(http.StatusOK, "Hello, World!")
}

func HandleMessages(e *echo.Echo) {
	g := e.Group("/messages")
	g.GET("/", CreateMessage)
}
