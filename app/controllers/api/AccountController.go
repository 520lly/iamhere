package controllers

import (
	"strconv"
	"time"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func HandleAccounts(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Accounts.Group
	g := e.Group(urlGroup)
	g.POST("/register", CreateNewAccount)
	g.PUT("/:id", UpdateAccount)
	g.POST("/login", ValidateAccount)
	g.GET("/login/", ValidateAccount)
	g.GET("/:id", GetAccounts)
	g.GET("", GetAccounts)
	g.DELETE("/:id", DeleteAccounts)
}

// Handlers
func CreateNewAccount(c echo.Context) error {
	var user User
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if err := DecodeBody(c, &user); err != nil {
		rsp.Code = RspBadRequest
		rsp.Reason = err.Error()
		RespondJ(c, RspBadRequest, rsp)
		return err
	}
	c.Logger().Debug(JsonToString(user))
	if err := HandleCreateNewUser(c, &user); err != nil {
		//rsp.Code = RspBadRequest
		//rsp.Reason = err.Error()
		//RespondJ(c, RspBadRequest, rsp)
		return err
	}
	return nil
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
	p := NewPath(c.Request().URL.Path)
	var user User
	if p.HasID() {
		if err := DecodeBody(c, &user); err != nil {
			rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
		}
		user.ID = StringToBson(p.GetID())
		if err := HandleUpdateUser(c, &user); err != nil {
			return err
		}
	} else {
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
	}
	return nil
}

func ValidateAccount(c echo.Context) error {
	//username could be associatedId/phonenumber/email query from URL.Query
	validatedPass := true
	lu := new(LoginUser)
	if err := c.Bind(lu); err != nil {
		//not Response immediately and check using URL Query
		lu.UserId = c.QueryParam("userid")
		lu.Password = c.QueryParam("password")
		usertypeS := c.QueryParam("usertype")
		if len(usertypeS) == 0 {
			rsp := &Response{RspBadRequest, ReasonMissingParam + "usertype", nil, 0}
			RespondJ(c, RspBadRequest, rsp)
			return NewError(ReasonMissingParam)
		}
		i, err := strconv.Atoi(usertypeS)
		if err != nil {
			rsp := &Response{RspBadRequest, ReasonFailureParam + "usertype", nil, 0}
			RespondJ(c, RspBadRequest, rsp)
			return NewError(ReasonMissingParam)
		}
		lu.UserType = i
		c.Logger().Debug("i = ", i)
		if i == UserType_Wechat {
			lu.JsCode = c.QueryParam("jscode")
			c.Logger().Debug("jscode = ", lu.JsCode)
			if len(lu.JsCode) == 0 {
				rsp := &Response{RspBadRequest, ReasonMissingParam + "jscode", nil, 0}
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonMissingParam)
			}
		}
	}
	c.Logger().Debug("LoginUser:", JsonToString(lu))

	if lu.UserType == UserType_Wechat {
		c.Logger().Debug("Requst Wechat Authentication")
		if err := RequstSessionAndOpenId(c, lu); err != nil {
			c.Logger().Debug(err.Error())
			validatedPass = false
		}
	} else {
		if err := LoginValidate(c, lu.UserId, lu.Password); err != nil {
			//LoginValidate failed
			validatedPass = false
			rsp := &Response{RspBadRequest, ReasonAuthFailed, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
			return NewError(ReasonAuthFailed)
		}
	}
	if validatedPass {
		// Create token
		token := CreateNewJWTToken()

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = lu.UserId
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(Config.ApiConfig.Secret))
		if err != nil {
			return err
		}
		rsp := &Response{RspOK, ReasonSuccess, map[string]string{
			"token": t,
		}, 0}
		RespondJ(c, RspOK, rsp)
	}
	return echo.ErrUnauthorized
}
