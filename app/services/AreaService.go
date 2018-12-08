package services

import (
	. "github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
)

func HandleCreateNewArea(c echo.Context, area *Area) error {
	if area == nil {
		return NewError("area is nil")
	}
	c.Logger().Debug("area: ", JsonToString(area))
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if len(BsonToString(area.ID)) == 0 {
		//it'a new area
		if len(area.Name) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "Name is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Name is empty")
		} else if len(area.Province) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "Province is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Province is empty")
		} else if len(area.City) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "City is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("City is empty")
		} else if len(area.District) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "District is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("District is empty")
		} else if len(area.Discription) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "Discription is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Discription is empty")
		} else if len(area.Address1) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "Address1 is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Address1 is empty")
		} else if !ValidateAreaCategory(area.Category) {
			rsp.Code = RspBadRequest
			rsp.Reason = "Category is not valide"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Category is valide")
		} else if !ValidateAreaType(area.Type) {
			c.Logger().Debug("type:", area.Type)
			rsp.Code = RspBadRequest
			rsp.Reason = "Type is not valide"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Type is not valid")
		} else if !CheckInRangefloat64(area.Latitude, LatitudeMinimum, LatitudeMaximum) {
			rsp.Code = RspBadRequest
			rsp.Reason = "Latitude is out of range"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Latitude is out of range")
		} else if !CheckInRangefloat64(area.Longitude, LongitudeMinimum, LongitudeMaximum) {
			rsp.Code = RspBadRequest
			rsp.Reason = "Longitude is out of range"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Longitude is out of range")
		}
		//else if !CheckInRangefloat64(area.Radius, RadiusMinimum, RadiusMaximum) {
		//rsp.Code = RspBadRequest
		//rsp.Reason = "Radius is out of range"
		//RespondJ(c, RspBadRequest, rsp)
		//return NewError("Radius is out of range")
		//}
		area.TimeStamp = CreateTimeStampUnix()
		area.Location.Coordinates = []float64{area.Longitude, area.Latitude}
		area.Location.Type = "Point"
		area.ID = CreateNewObjectId()
		if areaFound := FindAreaWithLocation(c, area.Longitude, area.Latitude); areaFound != nil {
			//It's stored Area or this area has overlap with other Areas
			rsp.Code = RspBadRequest
			rsp.Reason = ReasonDuplicate
			rsp.Data = areaFound
			rsp.Count = 1
			RespondJ(c, RspBadRequest, rsp)
			return NewError(ReasonDuplicate)
		} else {
			//It's a new Area
			if Insert(DBCAreas, area) {
				c.Logger().Debug("Insert DBCArea Success")
				if err := CreateGeoIndex(DBCAreas); err == nil {
					c.Logger().Debug("CreateGeoIndex Success")
					c.Logger().Debug("Response: ", JsonToString(rsp))
					RespondJ(c, RspOK, rsp)
					return nil
				} else {
					c.Logger().Debug("CreateGeoIndex failed! err:", err)
					rsp.Code = RspInternalServerError
					rsp.Reason = ReasonInsertFailure
					c.Logger().Debug("Response: ", JsonToString(rsp))
					RespondJ(c, RspInternalServerError, rsp)
					return nil
				}
			} else {
				//Insert area failed
				c.Logger().Debug("Insert DBCAreas  failed")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			}

		}
	} else {
		//[TODO]it'a stored area need to update area
	}
	return nil
}

func HandleDeleteAreas(c echo.Context, area *Area) error {
	c.Logger().Debug("Delete area ID :", area.ID)
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	err := DeleteAreaWithID(area.ID)
	if err != nil {
		c.Logger().Debug("Failed to delete area ID:", area.ID)
		rsp.Code = RspInternalServerError
		rsp.Reason = err.Error()
		RespondJ(c, RspInternalServerError, rsp)
		return nil
	}
	c.Logger().Debug("Succeed to delete messages ID:", area.ID)
	RespondJ(c, RspOK, rsp)
	return nil
}
func HandleUpdateArea(c echo.Context, area *Area) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	c.Logger().Debug("ToUpdateArea: ", JsonToString(area))
	changed := false
	if len(BsonToString(area.ID)) != 0 {
		var areaStored Area
		if err := Get(DBCAreas, area.ID, &areaStored); err != nil {
			c.Logger().Debug("Failed to find Area with area ID:", area.ID)
			rsp.Code = RspBadRequest
			rsp.Reason = err.Error()
			RespondJ(c, RspBadRequest, rsp)
			return nil
		}
		c.Logger().Debug("areaStored: ", JsonToString(areaStored))
		if len(area.Name) != 0 && area.Name != areaStored.Name {
			c.Logger().Debug("UpdateByIdField name: ", area.Name)
			if !UpdateByIdField(DBCAreas, area.ID, "name", area.Name) {
				//update name failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.Province) != 0 && area.Province != areaStored.Province {
			if !UpdateByIdField(DBCAreas, area.ID, "province", area.Province) {
				//update province failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.City) != 0 && area.City != areaStored.City {
			if !UpdateByIdField(DBCAreas, area.ID, "city", area.City) {
				//update city failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.District) != 0 && area.District != areaStored.District {
			if !UpdateByIdField(DBCAreas, area.ID, "district", area.District) {
				//update district failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.Discription) != 0 && area.Discription != areaStored.Discription {
			c.Logger().Debug("UpdateByIdField discription: ", area.Discription)
			if !UpdateByIdField(DBCAreas, area.ID, "discription", area.Discription) {
				//update discription failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.Address1) != 0 && area.Address1 != areaStored.Address1 {
			if !UpdateByIdField(DBCAreas, area.ID, "address1", area.Address1) {
				//update address1 failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(area.Address2) != 0 && area.Address2 != areaStored.Address2 {
			if !UpdateByIdField(DBCAreas, area.ID, "address2", area.Address2) {
				//update address2 failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if ValidateAreaCategory(area.Category) && area.Category != areaStored.Category {
			if !UpdateByIdField(DBCAreas, area.ID, "category", area.Category) {
				//update category failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if ValidateAreaType(area.Type) && area.Type != areaStored.Type {
			if !UpdateByIdField(DBCAreas, area.ID, "type", area.Type) {
				//update type failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if CheckInRangefloat64(area.Latitude, LatitudeMinimum, LatitudeMaximum) && area.Latitude != areaStored.Latitude {
			if !UpdateByIdField(DBCAreas, area.ID, "latitude", area.Latitude) {
				//update latitude failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			location := GeoJson{Type: areaStored.Location.Type, Coordinates: []float64{areaStored.Longitude, area.Latitude}}
			if !UpdateByIdField(DBCAreas, area.ID, "location", location) {
				//update location failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			CreateGeoIndex(DBCAreas)
			changed = true
		}
		if CheckInRangefloat64(area.Longitude, LongitudeMinimum, LongitudeMaximum) && area.Longitude != areaStored.Longitude {
			if !UpdateByIdField(DBCAreas, area.ID, "latitude", area.Latitude) {
				//update latitude failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			location := GeoJson{Type: areaStored.Location.Type, Coordinates: []float64{areaStored.Longitude, area.Latitude}}
			if !UpdateByIdField(DBCAreas, area.ID, "location", location) {
				//update location failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			CreateGeoIndex(DBCAreas)
			changed = true
		}
		//if CheckInRangefloat64(area.Radius, RadiusMinimum, RadiusMaximum) && area.Radius != areaStored.Radius {
		//if !UpdateByIdField(DBCAreas, area.ID, "radius", area.Radius) {
		////update radius failed
		//rsp.Code = RspBadRequest
		//rsp.Reason = ReasonOperationFailed
		//RespondJ(c, RspBadRequest, rsp)
		//return NewError(ReasonOperationFailed)
		//}
		//changed = true
		//}
		if changed {
			RespondJ(c, RspOK, rsp)
			return nil
		}
		rsp.Code = RspBadRequest
		rsp.Reason = ReasonDuplicate
		RespondJ(c, RspBadRequest, rsp)
		return NewError(ReasonDuplicate)
	} else {
		rsp.Code = RspBadRequest
		rsp.Reason = ReasonMissingParam
		RespondJ(c, RspBadRequest, rsp)
		return NewError(ReasonMissingParam)
	}
	return nil
}

func HandleGetAreas(c echo.Context, area *Area, debug bool) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if debug {
		//return all areas
		areas, err := FindAllArea()
		if err == nil {
			c.Logger().Debug("Found areas size:", len(areas))
			rsp.Data = areas
			rsp.Count = len(areas)
			RespondJ(c, RspOK, rsp)
			return nil
		}
		c.Logger().Debug("FindAllArea failed!!err:", err.Error())
	} else if areaFound := FindAreaWithLocation(c, area.Longitude, area.Latitude); areaFound != nil {
		c.Logger().Debug("Found areas :", JsonToString(areaFound))
		rsp.Data = areaFound
		rsp.Count = 1
		RespondJ(c, RspOK, rsp)
		return nil
	}
	rsp.Code = RspBadRequest
	rsp.Reason = ReasonNotFound
	RespondJ(c, RspBadRequest, rsp)
	return NewError("Not found any")
}

func FindAreaWithLocation(c echo.Context, lon float64, lat float64) (ret *Area) {
	areas, _ := FindAllArea()
	if areas != nil {
		for _, area := range areas {
			areaMatchs := FindAreas(lon, lat, area.Radius)
			if areaMatchs != nil {
				if len(areaMatchs) == 1 && areaMatchs[0].ID == area.ID {
					areaMatchJ := JsonToString(areaMatchs[0])
					c.Logger().Debug("found Area:", string(areaMatchJ))
					return areaMatchs[0]
				}
			}
		}
	}
	return nil
}
