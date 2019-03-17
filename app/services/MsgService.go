package services

import (
	. "github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
)

func HandleCreateNewMessage(c echo.Context, msg *Message) error {
	if msg == nil {
		return NewError("msg is nil")
	}
	c.Logger().Debug("msg longitude: ", msg.Longitude, " latitude: ", msg.Latitude, " altitude: ", msg.Altitude)
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if len(BsonToString(msg.ID)) == 0 {
		//it'a new message
		if len(msg.UserID) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "UserID is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("UserID is empty")
		} else if len(msg.Content) == 0 {
			rsp.Code = RspBadRequest
			rsp.Reason = "Content is empty"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Content is empty")
		} else if !CheckInRangefloat64(msg.Latitude, LatitudeMinimum, LatitudeMaximum) {
			rsp.Code = RspBadRequest
			rsp.Reason = "Latitude is out of range"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Latitude is out of range")
		} else if !CheckInRangefloat64(msg.Longitude, LongitudeMinimum, LongitudeMaximum) {
			rsp.Code = RspBadRequest
			rsp.Reason = "Longitude is out of range"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Longitude is out of range")
		} else if msg.ExpiryTime == 0 {
			msg.ExpiryTime = CreateTimeStampUnix()
		} else if msg.Color == 0 {
			msg.Color = Black
		} else if msg.Available {
			msg.Available = CanBeSeen
		}

		msg.TimeStamp = CreateTimeStampUnix()
		msg.Location.Coordinates = []float64{msg.Longitude, msg.Latitude}
		msg.Location.Type = "Point"
		msg.ID = CreateNewObjectId()
		if area := FindAreaWithLocation(c, msg.Longitude, msg.Latitude); area != nil {
			//this message belong to a specific area
			msg.AreaID = BsonToString(area.ID)
			if Insert(DBCAreaMessages, msg) {
				c.Logger().Debug("Insert ", msg.AreaID, " DBCAreaMessages Success")
				if err := CreateGeoIndex(DBCAreaMessages); err == nil {
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
				//Insert message failed
				c.Logger().Debug("Insert DBCAreaMessages  failed")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			}
		} else {
			//this message belong to ocean
			msg.AreaID = "Ocean"
			if Insert(DBCOceanMessages, msg) {
				c.Logger().Debug("Insert DBCOceanMessages Success")
				if err := CreateGeoIndex(DBCOceanMessages); err == nil {
					c.Logger().Debug("CreateGeoIndex Success")
					rsp := &Response{RspOK, ReasonSuccess, nil, 0}
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
				//Insert message failed
				c.Logger().Debug("Insert DBCOceanMessages failed")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			}
		}
	} else {
		//It'a stored message and we need to update it
		//[TODO] Missing updateing
		var found *Message = nil
		Get(DBCAreaMessages, msg.ID, found)
		if found == nil {
			//found in Ocean messages set
			Get(DBCOceanMessages, msg.ID, found)
			if found == nil {
				//not found this message
				c.Logger().Debug("Not found Message in DBCAreaMessages and DBCOceanMessages ")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			}
			//found in DBCOceanMessages
			//[TODO] Implemente update to DBCOceanMessages
		}
		//found in DBCAreaMessages
		//[TODO] Implemente update to DBCAreaMessages
	}
	return nil
}

func HandleGetMessages(c echo.Context, msg *Message, debug bool, page int, size int) error {
	c.Logger().Debug("msg:", msg, "  debug:", debug, "  size:", size, "  page:", page)
	var msgs []*Message
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if debug {
		//Found up to 10 ocean messages
		msgs = GetRandomMessages(DBCOceanMessages, size)
		if msgs == nil {
			c.Logger().Debug("Find up to ", size, " ocean messages failed")
			rsp := &Response{RspBadRequest, ReasonOperationFailed, nil, 0}
			RespondJ(c, RspInternalServerError, rsp)
			return nil
		} else {
			c.Logger().Debug("Failed find up to ", size, " ocean messages Success")
			rsp := &Response{RspOK, ReasonSuccess, &msgs, len(msgs)}
			RespondJ(c, RspOK, rsp)
			return nil
		}
	} else {
		if msg == nil {
			return NewError("msg is nil")
		}
		//Return a message with specific ID (Message ID comes first than location)
		if CheckBsonObjNotEmpty(msg.ID) {
			//var m Message
			if m := GetOneItemWithID(DBCAreaMessages, msg.ID, Message{}); m != nil {
				c.Logger().Debug("m:", m)
				//Found msg in Area collection and return
				rsp := &Response{RspOK, ReasonSuccess, &m, 1}
				RespondJ(c, RspOK, rsp)
			}
			if m := GetOneItemWithID(DBCOceanMessages, msg.ID, Message{}); m != nil {
				c.Logger().Debug("m:", msg)
				//Found msg in Ocean collection and return
				rsp := &Response{RspOK, ReasonSuccess, &m, 1}
				RespondJ(c, RspOK, rsp)
				return nil
			}
		} else if CheckStringNotEmpty(msg.AreaID) {
			if CheckStringNotEmpty(msg.AreaID) {
				if msgs, err := FindMsgsWith1Feild(msg.AreaID, "areaid", msg.AreaID, page, size); err == nil {
					c.Logger().Debug("Found msgs size: ", len(msgs))
					rsp.Data = msgs
					rsp.Count = len(msgs)
					RespondJ(c, RspOK, rsp)
					return nil
				}
			}
		} else if CheckStringNotEmpty(msg.UserID) {
			if msgs, err := GetUserMessages(msg.UserID, page, size); err == nil {
				c.Logger().Debug("Found msgs size: ", len(msgs))
				rsp.Data = msgs
				rsp.Count = len(msgs)
				RespondJ(c, RspOK, rsp)
				return nil
			}
		}

		c.Logger().Debug("msg longitude: ", msg.Longitude, " latitude: ", msg.Latitude)
		if !CheckInRangefloat64(msg.Latitude, LatitudeMinimum, LatitudeMaximum) {
			return NewError("Latitude is out of range")
		} else if !CheckInRangefloat64(msg.Longitude, LongitudeMinimum, LongitudeMaximum) {
			return NewError("Longitude is out of range")
		}

		if area := FindAreaWithLocation(c, msg.Longitude, msg.Latitude); area != nil {
			//this message belong to a specific area
			//Found up to random item limit area messages
			msgs = GetSpecifiedLocationMessages(DBCAreaMessages, msg.Longitude, msg.Latitude, area.Radius, Config.ApiConfig.RandomItemLimit)
			if msgs == nil {
				c.Logger().Debug("Find up to ", Config.ApiConfig.RandomItemLimit, " ocean messages failed")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			} else {
				c.Logger().Debug("Find up to ", Config.ApiConfig.RandomItemLimit, " ocean messages Success")
				RespondJ(c, RspOK, rsp)
				return nil
			}

		} else {
			//Found up to random item limit ocean messages
			msgs = GetRandomMessages(DBCOceanMessages, Config.ApiConfig.RandomItemLimit)
			if msgs == nil {
				c.Logger().Debug("Find up to ", Config.ApiConfig.RandomItemLimit, " ocean messages failed")
				rsp.Code = RspInternalServerError
				rsp.Reason = ReasonInsertFailure
				RespondJ(c, RspInternalServerError, rsp)
				return nil
			} else {
				c.Logger().Debug("Find up to ", Config.ApiConfig.RandomItemLimit, " ocean messages Success")
				rsp.Data = msgs
				rsp.Count = len(msgs)
				RespondJ(c, RspOK, rsp)
				return nil
			}
		}
	}
	//TODO Implemente paging
	return nil
}

func HandleUpdateMessage(c echo.Context, msg *Message) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	c.Logger().Debug("Update msg:", msg)
	changed := false
	if len(BsonToString(msg.ID)) != 0 {
		var msgStored Message
		if err := Get(DBCAreaMessages, msg.ID, &msgStored); err != nil {
			c.Logger().Debug("Failed to find msg with ID:", msg.ID)
			if err := Get(DBCOceanMessages, msg.ID, &msgStored); err != nil {
				if len(msg.LikeCount) != 0 && msg.LikeCount != msgStored.LikeCount {
					if !UpdateByIdField(DBCOceanMessages, msg.ID, "likecount", msg.LikeCount) {
						//update likecount failed
						rsp.Code = RspBadRequest
						rsp.Reason = ReasonOperationFailed
						RespondJ(c, RspBadRequest, rsp)
						return NewError(ReasonOperationFailed)
					}
					changed = true
				}
				if len(msg.Recommend) != 0 && msg.Recommend != msgStored.Recommend {
					if !UpdateByIdField(DBCOceanMessages, msg.ID, "recommend", msg.Recommend) {
						//update firstname failed
						rsp.Code = RspBadRequest
						rsp.Reason = ReasonOperationFailed
						RespondJ(c, RspBadRequest, rsp)
						return NewError(ReasonOperationFailed)
					}
					changed = true
				}
			} else {
				rsp.Code = RspBadRequest
				rsp.Reason = err.Error()
				RespondJ(c, RspBadRequest, rsp)
				return err
			}
		} else {
			if len(msg.LikeCount) != 0 && msg.LikeCount != msgStored.LikeCount {
				if !UpdateByIdField(DBCAreaMessages, msg.ID, "likecount", msg.LikeCount) {
					//update likecount failed
					rsp.Code = RspBadRequest
					rsp.Reason = ReasonOperationFailed
					RespondJ(c, RspBadRequest, rsp)
					return NewError(ReasonOperationFailed)
				}
				changed = true
			}
			if len(msg.Recommend) != 0 && msg.Recommend != msgStored.Recommend {
				if !UpdateByIdField(DBCAreaMessages, msg.ID, "recommend", msg.Recommend) {
					//update firstname failed
					rsp.Code = RspBadRequest
					rsp.Reason = ReasonOperationFailed
					RespondJ(c, RspBadRequest, rsp)
					return NewError(ReasonOperationFailed)
				}
				changed = true
			}
		}

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
}

func HandleDeleteMessages(c echo.Context, msg *Message) error {
	c.Logger().Debug("Delete message ID :", msg.ID)
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	err := DeleteMessageWithID(msg.ID)
	if err != nil {
		c.Logger().Debug("Failed to delete messages ID:", msg.ID)
		rsp.Code = RspInternalServerError
		rsp.Reason = err.Error()
		RespondJ(c, RspInternalServerError, rsp)
		return nil
	}
	c.Logger().Debug("Succeed to delete messages ID:", msg.ID)
	RespondJ(c, RspOK, rsp)
	return nil
}
