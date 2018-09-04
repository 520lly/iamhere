package controllers

import (
	"time"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

func HandleAccounts(e *echo.Echo) {
	g := e.Group("/accounts")
	g.POST("/register", CreateNewAccount)
	g.PUT("/:id", UpdateAccount)
	g.POST("/login", ValidateAccount)
	g.GET("/login/", ValidateAccount)
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
	lu := new(LoginUser)
	if err := c.Bind(lu); err != nil {
		//Missing failed
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
		return NewError(ReasonMissingParam)
	}
	c.Logger().Debug("LoginUser:", JsonToString(lu))

	if err := LoginValidate(c, lu.UserId, lu.Password); err != nil {
		//LoginValidate failed
		rsp := &Response{RspBadRequest, ReasonAuthFailed, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
		return NewError(ReasonAuthFailed)
	}
	// Create token
	token := CreateNewJWTToken()

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = lu.UserId
	claims["admin"] = false
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	rsp := &Response{RspOK, ReasonSuccess, map[string]string{
		"token": t,
	}, 0}
	RespondJ(c, RspOK, rsp)
	return echo.ErrUnauthorized
}
