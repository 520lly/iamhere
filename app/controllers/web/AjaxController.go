package controllers_web

import (
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
)

type Result struct {
	Ret    int         //结果状态码:0/1
	Reason string      //描述信息
	Data   interface{} //数据
}

type AjaxController struct {
}

/**
  处理登录请求
*/
func (controller *AjaxController) LoginAction(c echo.Context) {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	c.Response().Writer.Header().Set("content-type", "application/json") //设置数据格式为json
	err := c.Request().ParseForm()                                       //解析post请求
	if err != nil {
		rsp.Code = RspBadRequest
		rsp.Reason = ReasonFailureParam
		RespondJ(c, RspBadRequest, rsp)
		return
	}
	username := c.Request().FormValue("username")
	password := c.Request().FormValue("password")
	if username == "" || password == "" {
		rsp.Code = RspBadRequest
		rsp.Reason = ReasonMissingParam
		RespondJ(c, RspBadRequest, rsp)
		return
	}
	if err := LoginValidate(c, username, password); err != nil {
		rsp.Code = RspBadRequest
		rsp.Reason = ReasonNotFound
		RespondJ(c, RspBadRequest, rsp)
		return
	}

	//存入cookie

	cookie := http.Cookie{Name: "username", Value: username, Path: "/"}
	http.SetCookie(c.Response().Writer, &cookie)
	//OutputJson(c.Response().Writer, 1, "操作成功", nil)
	RespondJ(c, RspOK, rsp)
	c.Logger().Debug("LoginAction Success")
	return
}

/**
  输出json
*/
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := &Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}
	w.Write(b)
}

/**
  得到所有用户
*/
func (controller *AjaxController) GetUsersAction(w http.ResponseWriter, r *http.Request) {
	//func (controller *AjaxController) GetUsersAction(c echo.Context) {
	//c.Response().Writer.Header().Set("content-type", "application/json") //设置数据格式为json
	w.Header().Set("content-type", "application/json") //设置数据格式为json
	err := r.ParseForm()                               //解析post请求
	//err := c.Request().ParseForm()                                       //解析post请求
	if err != nil {
		//OutputJson(c.Response().Writer, 0, "参数错误", nil)
		OutputJson(w, 0, "参数错误", nil)
		return
	}

	//db := NewDataMng()
	//if err := db.Connect(); err != nil {
	//    log.Println(err)
	//    OutputJson(w, 0, "数据库操作失败", nil)
	//    return
	//}
	//defer db.Close()
	//rows, res, err := db.Query("select * from user where username='%s'", "LCore")
	//if err != nil {
	//    log.Println(err)
	//    OutputJson(w, 0, "数据库操作失败", nil)
	//    return
	//}
	//OutputJson(c.Response().Writer, 1, "操作成功", "root")
	OutputJson(w, 1, "操作成功", "root")
	return
}
