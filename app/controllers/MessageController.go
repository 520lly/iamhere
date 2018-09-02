package controllers

import (
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/labstack/echo"
)

func HandleMessages(e *echo.Echo) {
	g := e.Group("/messages")
	g.POST("/", CreateNewMessage)
	g.GET("/", GetMessages)
	g.DELETE("/:id", DeleteMessages)
}

// Handlers
func CreateNewMessage(c echo.Context) error {
	var msg Message
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if err := DecodeBody(c, &msg); err != nil {
		rsp.Code = RspBadRequest
		rsp.Reason = err.Error()
		RespondJ(c, RspBadRequest, rsp)
		return err
	}
	c.Logger().Debug(JsonToString(msg))
	if err := HandleCreateNewMessage(c, &msg); err != nil {
		return err
	}
	return nil
}

//Handler for GetMessages
func GetMessages(c echo.Context) error {
	var msg Message
	var debugF bool = false
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if err := DecodeBody(c, &msg); err != nil {
		//not Response immediately and check using URL Query
		debug := c.QueryParam("debug")
		c.Logger().Debug("debug:", debug)
		if len(debug) == 0 {
			if msg.Longitude, err = ConvertString2Float64(c.QueryParam("longitude")); err != nil {
				rsp.Code = RspBadRequest
				rsp.Reason = err.Error()
				RespondJ(c, RspBadRequest, rsp)
				return err
			}
			if msg.Latitude, err = ConvertString2Float64(c.QueryParam("latitude")); err != nil {
				rsp.Code = RspBadRequest
				rsp.Reason = err.Error()
				RespondJ(c, RspBadRequest, rsp)
				return err
			}
			c.Logger().Debug(JsonToString(msg))
		} else {
			debugF = true
		}
	}
	if err := HandleGetMessages(c, &msg, debugF); err != nil {
		return err
	}
	return nil
}

//Handler for Delete Messages
func DeleteMessages(c echo.Context) error {
	c.Logger().Debug("DeleteMessages called")
	p := NewPath(c.Request().URL.Path)
	if p.HasID() {
		msg := &Message{ID: StringToBson(p.GetID())}
		if err := HandleDeleteMessages(c, msg); err != nil {
			return err
		} else {
			rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
		}
	}
	return nil
}
