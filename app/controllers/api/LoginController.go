package controllers

import (
	"strconv"
	"time"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

func HandleLogin(e *echo.Echo) {
	urlGroup := Config.ApiConfig.Prefix + Config.ApiConfig.Version + Config.ApiConfig.Login.Group
	g := e.Group(urlGroup)
	//g.Use(middleware.JWT(GetJWTSecretCode()))
	g.POST("/register", CreateNewAccount)
	g.POST("", ValidateAccount)
	g.GET("", ValidateAccount)
}

// Handlers
func CreateNewAccount(c echo.Context) error {
	var user User
	if err := DecodeBody(c, &user); err != nil {
		return err
	}
	c.Logger().Debug(JsonToString(user))
	if err := HandleCreateNewUser(c, &user); err != nil {
		return err
	}
	return nil
}

func ValidateAccount(c echo.Context) error {
	//username could be associatedId/phonenumber/email query from URL.Query
	validatedPass := true
	lu := new(LoginUser)
	user := new(User)
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

	//if usertype is Wechat, then Requst to Authentication and then Create a new user/ValidateAccount if successfully
	if lu.UserType == UserType_Wechat {
		c.Logger().Debug("Requst Wechat Authentication")
		if err, woi := RequstSessionAndOpenId(c, lu); err != nil {
			c.Logger().Debug("RequstSessionAndOpenId Err: ", err.Error())
			validatedPass = false
			rsp := &Response{RspBadRequest, err.Error(), nil, 0}
			RespondJ(c, RspBadRequest, rsp)
		} else {
			//handle expired_in data from Wechat server
			lu.UserId = woi.OpenId
			if user, err = LoginValidate(c, lu.UserId, lu.Password); err != nil {
				//Not a registered user
				user := User{AssociatedId: woi.OpenId, Password: lu.Password}
				if err := HandleCreateNewUser(c, &user); err != nil {
					c.Logger().Debug("HandleCreateNewUser Err: ", err.Error())
					validatedPass = false
					rsp := &Response{RspBadRequest, ReasonAuthFailed, nil, 0}
					RespondJ(c, RspBadRequest, rsp)
					return err
				}
				//Create User in iamhere with Wechat appid successfully then return token
			} else {
				//This is a registered user and return token
			}
		}
	} else {
		var err error
		if user, err = LoginValidate(c, lu.UserId, lu.Password); err != nil {
			validatedPass = false
			rsp := &Response{RspBadRequest, ReasonAuthFailed, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
			return NewError(ReasonAuthFailed)
		} else {
			//This is a registered user and return token
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
			"token":  t,
			"userid": BsonToString(user.ID),
		}, 0}
		RespondJ(c, RspOK, rsp)
		return nil
	}
	return echo.ErrUnauthorized
}
