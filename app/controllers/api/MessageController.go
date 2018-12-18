package controllers

import (
	"strconv"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func HandleMessages(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Messages.Group
	g := e.Group(urlGroup)
	g.Use(middleware.JWT(GetJWTSecretCode()))
	g.POST("", CreateNewMessage)
	g.GET("/:id", GetMessages)
	g.PUT("/:id", UpdateMessages)
	g.GET("", GetMessages)
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
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	var msg Message
	var debugF bool = false
	if err := DecodeBody(c, &msg); err != nil {
		//not Response immediately and check using URL Query
		debug := c.QueryParam("debug")
		c.Logger().Debug("debug:", debug)
		if len(debug) != 0 {
			debugF = true
			c.Logger().Debug("debugF:", debugF)
		} else {
			longitude := c.QueryParam("longitude")
			if msg.Longitude, err = ConvertString2Float64(longitude); err == nil {
				msg.Longitude = longitude
				c.Logger().Debug("longtitude:", longtitude)
			}
			latitude := c.QueryParam("latitude")
			if msg.Latitude, err = ConvertString2Float64(latitude); err == nil {
				msg.Latitude = latitude
				c.Logger().Debug("latitude:", latitude)
			}
			areaid := c.QueryParam("areaid")
			if len(areaid) != 0 {
				c.Logger().Debug("areaid:", areaid)
				msg.AreaID = areaid
			}
			userid := c.QueryParam("userid")
			if len(userid) != 0 {
				c.Logger().Debug("userid:", userid)
				msg.UserID = userid
			}
			c.Logger().Debug("msg:  ", JsonToString(msg))
		}
	}
	p := NewPath(c.Request().URL.Path)
	c.Logger().Debug("p.HasID:", p.HasID(), "   p.GetID: ", p.GetID())
	if p.HasID() {
		msg.ID = ConvertString2BsonObjectId(p.GetID())
	}
	if err := HandleGetMessages(c, &msg, debugF); err != nil {
		return err
	}
	return nil
}

func UpdateMessages(c echo.Context) error {
	var msg Message
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if err := c.Bind(msg); err != nil {
		//not Response immediately and check using URL Query
		p := NewPath(c.Request().URL.Path)
		if p.HasID() {
			likecount := c.QueryParam("likecount")
			if len(likecount) != 0 {
				lc, err := strconv.ParseInt(likecount, 0, 32)
				if err != nil {
					rsp.Code = RspBadRequest
					rsp.Reason = err.Error()
					RespondJ(c, RspBadRequest, rsp)
					return NewError(ReasonMissingParam)
				}
				msg.LikeCount = likecount
				c.Logger().Debug("likecount=", lc)
			}
			recommend := c.QueryParam("recommend")
			c.Logger().Debug("recommend=", recommend)
			if len(recommend) != 0 {
				rc, err := strconv.ParseBool(recommend)
				if err != nil {
					rsp.Code = RspBadRequest
					rsp.Reason = err.Error()
					RespondJ(c, RspBadRequest, rsp)
					return NewError(ReasonMissingParam)
				}
				msg.Recommend = recommend
				c.Logger().Debug("recommend=", rc)
			}
		} else {
			c.Logger().Debug("Not specific Message ID")
			rsp.Code = RspBadRequest
			rsp.Reason = "Not specific Message ID"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Not specific Message ID")
		}
	}

	c.Logger().Debug(JsonToString(msg))
	if err := HandleUpdateMessage(c, &msg); err != nil {
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
