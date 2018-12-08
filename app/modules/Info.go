package modules

import (
	"net/http"
)

const (
	LatitudeMinimum float64 = -90 //Latitude Minimum
	LatitudeMaximum float64 = 90  //Latitude Maximum
)

const (
	LongitudeMinimum float64 = -180 //Longitude Minimum
	LongitudeMaximum float64 = 180  //Longitude Maximum
)

const (
	AltitudeMinium  float64 = -10000 //Altitude Minimum
	AltitudeMaximum float64 = 10000  //Altitude Maximum
)

const (
	RadiusMinimum float64 = 10    // Radius Minimum meters
	RadiusMaximum float64 = 10000 // Radius Maximum
)

const (
	CategoryMinimum int = -1               //Category Minimum
	CategoryFixed   int = 1                //defined by system for fixed areas
	CategoryOcean   int = 2                //defined by system for ocean
	CategoryIsland  int = 4                //defined by system for island
	CategoryCloud   int = 8                //defined by system for cloud
	CategoryUser    int = 16               //defined by user
	CategoryMaximum int = CategoryUser + 1 //Category Maximum
)

const (
	TypeMinimum    int = -1                 //Area type Minimum
	TypeNatural    int = 0                  //Area type is Natural
	TypeHistorical int = 1                  //Area type is Historical
	TypeHumanities int = 2                  //Area type is Humanities
	TypeMaximum    int = TypeHumanities + 1 //Area type Maximum
)

const (
	RspOK                  int = http.StatusOK
	RspBadRequest          int = http.StatusBadRequest
	RspInternalServerError int = http.StatusInternalServerError
)
const (
	ReasonSuccess         string = "Success "
	ReasonFailureParam    string = "Wrong parameter: "
	ReasonMissingParam    string = "Missing parameter: "
	ReasonFailureAPIKey   string = "Wrong APIKey: "
	ReasonFailueGeneral   string = "Failure in general: "
	ReasonDuplicate       string = "Parameter duplicated: "
	ReasonInsertFailure   string = "Insert failed: "
	ReasonWrongPw         string = "Wrong Password:  "
	ReasonOperationFailed string = "Operation Failure:  "
	ReasonNotFound        string = "Not Found: "
	ReasonAlreadyExist    string = "Already Existed: "
	ReasonAuthFailed      string = "Authentication failed:  "
)
const (
	RandomItemLimit int = 10
)

const (
	UserType_Wechat      int = 0                  //The user account type is wechat
	UserType_PhoneNumber int = 1                  //The user account type is phonenumber
	UserType_Email       int = 2                  //The user account type is email
	UserType_Other       int = UserType_Email + 1 //The user account type is other
)

type GeoJson struct {
	Type        string    `josn:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Response struct {
	Code   int         `json:"code"`
	Reason string      `json:"reasone"`
	Data   interface{} `json:"data"`
	Count  int         `json:"count"`
}

// Login User
type LoginUser struct {
	UserId   string `json:"userid" form:"userid" query:"userid"`       //for Wechat user, userid is weapp appid
	Password string `json:"password" form:"password" query:"password"` // for wechat user, password is secret
	JsCode   string `json:"jscode" form:"jscode" query:"jscode"`       // for wechat user, the value matters or make no sense
	UserType int    `json:"usertype" form:"usertype" query:"usertype"` // for wechat user, the usetype should be UserType_Wechat(0)
}

type WechatOpenId struct {
	OpenId     string `json:"openid" from:"openid" query:"openid"`
	SessionKey string `json:"session_key" from:"session_key" query:"session_key"`
	ExpiresIn  int    `json:"expired_in" from:"expired_in" query:"expired_in"`
}
