package controllers

import (
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	. "github.com/520lly/iamhere/app/services"
	"github.com/labstack/echo"
)

func HandleAreas(e *echo.Echo) {
	g := e.Group("/areas")
	g.POST("/", CreateNewArea)
	g.POST("/:id", UpdateArea)
	g.GET("/", GetAreas)
	g.DELETE("/:id", DeleteAreas)
}

// Handlers
func CreateNewArea(c echo.Context) error {
	var area Area
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if err := DecodeBody(c, &area); err != nil {
		rsp.Code = RspBadRequest
		rsp.Reason = err.Error()
		RespondJ(c, RspBadRequest, rsp)
		return err
	}
	c.Logger().Debug(JsonToString(area))
	if err := HandleCreateNewArea(c, &area); err != nil {
		//rsp.Code = RspBadRequest
		//rsp.Reason = err.Error()
		//RespondJ(c, RspBadRequest, rsp)
		return err
	}
	return nil
}

//Handler for GetAreas
func GetAreas(c echo.Context) error {
	var area Area
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	var debugF bool = false
	if err := DecodeBody(c, &area); err != nil {
		//not Response immediately and check using URL Query
		debug := c.QueryParam("debug")
		c.Logger().Debug("debug:", debug)
		if len(debug) == 0 {
			if area.Longitude, err = ConvertString2Float64(c.QueryParam("longitude")); err != nil {
				c.Logger().Debug("err:", err.Error())
				rsp.Code = RspBadRequest
				rsp.Reason = err.Error()
				RespondJ(c, RspBadRequest, rsp)
				return err
			}
			if area.Latitude, err = ConvertString2Float64(c.QueryParam("latitude")); err != nil {
				c.Logger().Debug("err:", err.Error())
				rsp.Code = RspBadRequest
				rsp.Reason = err.Error()
				RespondJ(c, RspBadRequest, rsp)
				return err
			}
			c.Logger().Debug(JsonToString(area))
		} else {
			//Get all areas for debugging purpose
			debugF = true
		}
	}
	if err := HandleGetAreas(c, &area, debugF); err != nil {
		c.Logger().Debug(err.Error())
		//rsp.Code = RspBadRequest
		//rsp.Reason = err.Error()
		//RespondJ(c, RspBadRequest, rsp)
		return err
	}
	return nil
}

//Handler for Delete Areas
func DeleteAreas(c echo.Context) error {
	p := NewPath(c.Request().URL.Path)
	if p.HasID() {
		area := &Area{ID: StringToBson(p.GetID())}
		if err := HandleDeleteAreas(c, area); err != nil {
			return err
		}
	} else {
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
	}
	return nil
}

//Handler for updating area
func UpdateArea(c echo.Context) error {
	p := NewPath(c.Request().URL.Path)
	var area Area
	if p.HasID() {
		if err := DecodeBody(c, &area); err != nil {
			rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
			RespondJ(c, RspBadRequest, rsp)
		}
		area.ID = StringToBson(p.GetID())
		if err := HandleUpdateArea(c, &area); err != nil {
			return err
		}
	} else {
		rsp := &Response{RspBadRequest, ReasonMissingParam, nil, 0}
		RespondJ(c, RspBadRequest, rsp)
	}
	return nil
}
