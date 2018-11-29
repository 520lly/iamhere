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

func RequstSessionAndOpenId(c echo.Context, lu *LoginUser) error {
	fullUrl := Config.ExternalUrl.Wechat.Url + "?appid=" + lu.UserId + "&secret=" + lu.Password + "js_code=" + lu.JsCode + "grant_type" + Config.ExternalUrl.Wechat.GrantType
	c.Logger().Debug("fullUrl is --->", fullUrl)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(fullUrl)
	if err != nil {
		// handle error
		return err
	}
	//return Success then parse response.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return err
	}
	var wechatOpenId WechatOpenId
	json.Unmarshal(body, &wechatOpenId)
	c.Logger().Debug("wechatOpenId", wechatOpenId)
	if wechatOpenId.OpenId == "" || wechatOpenId.SessionKey == "" {
		// handle error
		return errors.New("Error:Wrong jscode/secret/appid")
	}
	return nil
}
