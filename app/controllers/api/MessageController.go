package controllers

import (
	//"fmt"
	"strconv"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/dgrijalva/jwt-go"
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

	user := c.Get("user").(*jwt.Token)
	if user != nil {
		claims := user.Claims.(jwt.MapClaims)
		c.Logger().Debug("msg.UserID :", claims["name"])
		msg.UserID = claims["name"].(string)
	}
	c.Logger().Debug(JsonToString(msg))
	if err := HandleCreateNewMessage(c, &msg); err != nil {
		return err
	}
	return nil
}

//Handler for GetMessages
func GetMessages(c echo.Context) error {
	//rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	var msg Message
	var debugF bool = false
	var sizeLimit = Config.ApiConfig.RandomItemLimit
	var page = 0
	var err error
	//if err := DecodeBody(c, &msg); err != nil {
	//not Response immediately and check using URL Query
	debug := c.QueryParam("debug")
	c.Logger().Debug("debug:", debug)
	if len(debug) != 0 {
		debugF = true
		c.Logger().Debug("debugF:", debugF)
	} else {
		long := c.QueryParam("longitude")
		c.Logger().Debug("longitude:", long)
		if CheckStringNotEmpty(long) {
			if msg.Longitude, err = ConvertString2Float64(long); err == nil {
				c.Logger().Debug("longitude:", msg.Longitude)
			}
		} else {
			msg.Longitude = LongitudeMinimum
		}
		lat := c.QueryParam("latitude")
		if CheckStringNotEmpty(lat) {
			if msg.Latitude, err = ConvertString2Float64(lat); err == nil {
				c.Logger().Debug("latitude:", msg.Latitude)
			}
		} else {
			msg.Latitude = LatitudeMinimum
		}
		areaid := c.QueryParam("areaid")
		if len(areaid) != 0 {
			c.Logger().Debug("areaid:", areaid)
			msg.AreaID = areaid
		} else {
         user := c.Get("user").(*jwt.Token)
         if user != nil {
            claims := user.Claims.(jwt.MapClaims)
            c.Logger().Debug("msg.UserID :", claims["name"])
            msg.UserID = claims["name"].(string)
         }
      }
		sizeLimit, _ = strconv.Atoi(c.QueryParam("size"))
		if !CheckSizeLimitValidate(sizeLimit) {
			sizeLimit = Config.ApiConfig.RandomItemLimit
		}
		page, _ = strconv.Atoi(c.QueryParam("page"))
	}
	//}
	id := c.Param("id")
	if CheckStringNotEmpty(id) {
		msg.ID = StringToBson(id)
	}
	c.Logger().Debug("msg:  ", JsonToString(msg))
	if err := HandleGetMessages(c, &msg, debugF, page, sizeLimit); err != nil {
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
