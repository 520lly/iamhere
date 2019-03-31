package messages

import (
	. "github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"

	"github.com/labstack/echo"
)

var (
   //0x0000000
   GetMessagesChoice int32 = 0
)

const (
   ChoiceFlagSelected int32 = 1
   ChoiceDebugSelected int32 = 1
   ChoiceMsgIDSelected int32 = 2
   ChoiceAreaIDSelected int32 = 3
   ChoiceUserIDSelected int32 = 4
   ChoiceGeoValidShife int32 = 5
)


func HandleGetMessages(c echo.Context, msg *Message, debug bool, page int, size int) error {
   c.Logger().Debug("msg:", msg, "  debug:", debug, "  size:", size, "  page:", page)
   //msgs found to Response
   var msgs []*Message 
   var retErr error = nil
   rsp := &Response{RspOK, ReasonSuccess, nil, 0}

   if debug {
      GetMessagesChoice | (ChoiceFlagSelected << ChoiceDebugSelected)
   } 
   if msg == nil {
      retErr = NewError("params error: msg is nil")
   } else {
      //check if specified a message ID
      if CheckBsonObjNotEmpty(msg.ID) {
         GetMessagesChoice | (ChoiceFlagSelected << ChoiceMsgIDSelected)
      }
      //check if specified a area ID
      if CheckBsonObjNotEmpty(msg.AreaID) {
         GetMessagesChoice | (ChoiceFlagSelected << ChoiceAreaIDSelected)
      }
      //check if specified a user ID
      if CheckStringNotEmpty(msg.UserID) {
         GetMessagesChoice | (ChoiceFlagSelected << ChoiceUserIDSelected)
      }
      //verify GEO information 
      c.Logger().Debug("msg longitude: ", msg.Longitude, " latitude: ", msg.Latitude)
      if CheckInRangefloat64(msg.Latitude, LatitudeMinimum, LatitudeMaximum) && CheckInRangefloat64(msg.Longitude, LongitudeMinimum, LongitudeMaximum) {
         GetMessagesChoice | (ChoiceFlagSelected << ChoiceGeoValidShife)
      }
   }

   c.Logger().Debug("GetMessagesChoice: [", GetMessagesChoice, "]")
   switch GetMessagesChoice {
   case ChoiceMsgIDSelected:
      retErr = getItemWithMsgID(msg,msgs)
   case ChoiceUserIDSelected:
      retErr = getItemsWithUserID(page, size, msg, msgs)
   case ChoiceUserIDSelected + ChoiceAreaIDSelected:
      retErr = getItemsWithUserAndAreaID(page, size, msg, msgs)
   case ChoiceAreaIDSelected:
      retErr = getItemsWithAreaID(page, size, msg,msgs)
   case ChoiceMsgIDSelected:
      retErr = getItemWithMsgID(msg,msgs)
   case ChoiceMsgIDSelected:
      retErr = getItemWithMsgID(msg,msgs)
   case ChoiceMsgIDSelected:
      retErr = getItemWithMsgID(msg,msgs)
   }
   //getRandomItems(c, size, msgs)
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
   //TODO Implemente paging
   return nil
}

func getItemWithMsgID(msg *Message, msgs []*Message/*Messages found and return*/) error {
   err := nil
   //Try to get specified message within area message list
   if m := GetOneItemWithID(DBCAreaMessages, msg.ID, Message{}); m != nil {
      msgs.append(msgs,m...)
   } else {
      //Try to get specified message within ocean message list
      if m := GetOneItemWithID(DBCOceanMessages, msg.ID, Message{}); m != nil {
         msgs.append(msgs,m...)
      } else {
         err = NewError("Operation failed: Failed to get specified msg from area msgs list!")
      }
   }
   return err
}

func getItemsWithUserID(page int, size int, msg *Message, msgs []*Message/*Messages found and return*/) error {
   err := nil
   if msgs, err = GetUserMessages(msg.UserID, page, size); err != nil {
      err = NewError("Operation failed: Failed to get msgs with specified user ID from area/ocean msgs list!")
   }
   return err
}

func getItemsWithUserAndAreaID(page int, size int, msg *Message, msgs []*Message/*Messages found and return*/) error {
   err := nil
   m := make(map[string]string)
   m["key1"] = "userid"
   m["value1"] = msg.UserID
   m["key2"] = "areaid"
   m["value2"] = msg.AreaID
   m["key3"] = "limitaccess"
   m["value3"] = false
   if msgs, err = FindMsgsWith2Feild(msg.AreaID, "areaid", msg.AreaID, page, size); err != nil {
   if msgs, err = FindMsgsWith2Feild(msg.AreaID, "areaid", msg.AreaID, page, size); err != nil {
      err = NewError("Operation failed: Failed to get msgs with specified user ID from area/ocean msgs list!")
   }
   return err
}

func getItemsWithAreaID(page int, size int, msg *Message, msgs []*Message/*Messages found and return*/) error {
   err := nil
   if msgs, err = FindMsgsWith1Feild(msg.AreaID, "areaid", msg.AreaID, page, size); err != nil {
      err = NewError("Operation failed: Failed to get msgs with specified area ID from msgs list!")
   }
   return err
}



func getRandomItems(size int, msgs []*Message /*Messages found and return*/) error {
   msgs = GetRandomMessages(DBCOceanMessages, size)
   if msgs == nil {
      return NewError("Operation failed: Failed to get random messages from Ocean list!")
   } else {
      return nil
   }
}
