package controllers

import (
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func HandleAccounts(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Accounts.Group
	g := e.Group(urlGroup)
	g.Use(middleware.JWT(GetJWTSecretCode()))
	g.PUT("/update", UpdateAccount)
	g.GET("/:id", GetAccounts)
	g.GET("", GetAccounts)
	g.DELETE("/:id", DeleteAccounts)
}

//Handler for GetAccounts
func GetAccounts(c echo.Context) error {
	var user User
	var debugF bool = false
	if err := DecodeBody(c, &user); err != nil {
		//not Response immediately and check using URL Query
		debug := c.QueryParam("debug")
		c.Logger().Debug("debug:", debug)
		if len(debug) == 0 {
			p := NewPath(c.Request().URL.Path)
			if p.HasID() {
				user.ID = StringToBson(p.GetID())
			} else {
				rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonMissingParam)
			}
			c.Logger().Debug(JsonToString(user))
		} else {
			//Get all accounts for debugging purpose
			debugF = true
		}
	}
	if err := HandleGetUsers(c, &user, debugF); err != nil {
		c.Logger().Debug(err.Error())
		//rsp.Code = RspBadRequest
		//rsp.Reason = err.Error()
		//RespondJ(c, RspBadRequest, rsp)
		return err
	}
	return nil
}

//Handler for Delete Accounts
func DeleteAccounts(c echo.Context) error {
	p := NewPath(c.Request().URL.Path)
	if p.HasID() {
		user := &User{ID: StringToBson(p.GetID())}
		if err := HandleDeleteUsers(c, user); err != nil {
			return err
		}
	} else {
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
	}
	return nil
}

//Handler for updating user
func UpdateAccount(c echo.Context) error {
	var user User
	if err := DecodeBody(c, &user); err == nil {
		method := c.QueryParam("method")
		c.Logger().Debug("method: ", method)
		id := c.Param("id")
		if CheckStringNotEmpty(id) {
			c.Logger().Debug("id: ", id)
			user.ID = StringToBson(id)
			if err := HandleUpdateUser(c, &user, method); err != nil {
				return err
			}
		} else if userToken := c.Get("user").(*jwt.Token); userToken != nil {
			c.Logger().Debug("id: ", id)
			//:= c.Get("user").(*jwt.Token)
			c.Logger().Debug("userToken: ", userToken)
			if userToken != nil {
				claims := userToken.Claims.(jwt.MapClaims)
				c.Logger().Debug("UserID :", claims["name"])
				if u, err := GetAccountIDViaUserID(c, claims["name"].(string)); err == nil {
					c.Logger().Debug("user found: ", u)
					user.ID = u.ID
					c.Logger().Debug("ID: ", user.ID)
					if err := HandleUpdateUser(c, &user, method); err != nil {
						return err
					}
				} else {
					return NewError("not found")
				}
			}
		} else {
			rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
		}
	} else {
		c.Logger().Debug("err: ", err.Error())
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
	}
	return nil
}
