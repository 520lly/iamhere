package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	//"strconv"
	"errors"
	"time"
	//. "github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
)

func RequstSessionAndOpenId(c echo.Context, lu *LoginUser) (error, WechatOpenId) {
	fullUrl := Config.ExternalUrl.Wechat.Url + "?appid=" + lu.UserId + "&secret=" + lu.Password + "&js_code=" + lu.JsCode + "&grant_type=" + Config.ExternalUrl.Wechat.GrantType
	c.Logger().Debug("fullUrl is --->", fullUrl)
	var wechatOpenId WechatOpenId
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullUrl)
	if err != nil {
		// handle error
		return err, wechatOpenId
	}
	//return Success then parse response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return err, wechatOpenId
	}
	json.Unmarshal(body, &wechatOpenId)
	c.Logger().Debug("resp.Body=", wechatOpenId)
	c.Logger().Debug("wechatOpenId.OpenId=", wechatOpenId.OpenId)
	//if wechatOpenId.OpenId == "" || wechatOpenId.SessionKey == "" {
	if wechatOpenId.OpenId == "" {
		// handle error
		return errors.New("Error:Wrong jscode/secret/appid"), wechatOpenId
	}
	return nil, wechatOpenId
}
