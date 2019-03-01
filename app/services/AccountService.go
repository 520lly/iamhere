package services

import (
	. "github.com/520lly/iamhere/app/db"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
)

func HandleCreateNewUser(c echo.Context, user *User) error {
	if user == nil {
		return NewError("user is nil")
	}
	c.Logger().Debug("user: ", JsonToString(user))
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if len(BsonToString(user.ID)) == 0 {
		//check user duplicated
		found, err := isUserRegistered(c, user)
		c.Logger().Debug("found:", found)
		if found && err == nil {
			rsp.Code = RspBadRequest
			rsp.Reason = ReasonAlreadyExist
			RespondJ(c, RspBadRequest, rsp)
			//return NewError(ReasonAlreadyExist)
			return nil
		}

		//it'a new user
		if err := registerUser(c, user); err != nil {
			rsp.Code = RspBadRequest
			rsp.Reason = err.Error()
			RespondJ(c, RspBadRequest, rsp)
			return NewError(err.Error())
		}
	} else {
		//[TODO]it'a stored user need to update user
	}
	return nil
}

func registerUser(c echo.Context, user *User) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if len(user.AssociatedId) == 0 && len(user.PhoneNumber) == 0 && len(user.Email) == 0 {
		return NewError("associatedId/phonenumber/email must be valid at least one")
	} else if len(user.NickName) == 0 {
		user.NickName = CreateRandomNickname()
	} else if len(user.Password) == 0 {
		return NewError("Password is empty")
	}

	user.TimeStamp = CreateTimeStampUnix()
	user.ID = CreateNewObjectId()
	if Insert(DBCAccounts, user) {
		c.Logger().Debug("Insert DBCAccounts Success")
		RespondJ(c, RspOK, rsp)
		return nil
	} else {
		c.Logger().Debug("Insert DBCAccounts Failed")
		return NewError(ReasonInsertFailure)
	}
	return nil
}

func isUserRegistered(c echo.Context, user *User) (bool, error) {
	var usersFound []*User
	m := make(map[string]string)
	if len(user.AssociatedId) != 0 {
		m["key1"] = "associatedId"
		m["value1"] = user.AssociatedId
	}
	if len(user.PhoneNumber) != 0 {
		m["key2"] = "phonenumber"
		m["value2"] = user.PhoneNumber
	}
	if len(user.Email) != 0 {
		m["key3"] = "email"
		m["value3"] = user.Email
	}
	c.Logger().Debug("map side:", len(m))
	if usersFound, err := FindUsersWithFeild(DBCAccounts, m); err == nil {
		if len(usersFound) != 0 {
			c.Logger().Debug("Found users size: ", len(usersFound))
			return true, err
		}
	}
	c.Logger().Debug("Found users size: ", len(usersFound))
	return false, nil
}

func HandleDeleteUsers(c echo.Context, user *User) error {
	c.Logger().Debug("Delete user ID :", user.ID)
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	err := DeleteAccountWithID(user.ID)
	if err != nil {
		c.Logger().Debug("Failed to delete user ID:", user.ID)
		rsp.Code = RspInternalServerError
		rsp.Reason = err.Error()
		RespondJ(c, RspInternalServerError, rsp)
		return nil
	}
	c.Logger().Debug("Succeed to delete messages ID:", user.ID)
	RespondJ(c, RspOK, rsp)
	return nil
}

func HandleUpdateUser(c echo.Context, user *User, method string) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	changed := false
	if len(BsonToString(user.ID)) != 0 {
		var userStored User
		if userStored := GetOneItemWithID(DBCAccounts, user.ID, User{}); userStored == nil {
			c.Logger().Debug("Failed to find user with user ID:", user.ID)
			rsp.Code = RspBadRequest
			rsp.Reason = "Not found"
			RespondJ(c, RspBadRequest, rsp)
			return NewError("Not found")
		}
		if len(user.NickName) != 0 && user.NickName != userStored.NickName {
			if !UpdateByIdField(DBCAccounts, user.ID, "nickname", user.NickName) {
				//update nickname failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.FirstName) != 0 && user.FirstName != userStored.FirstName {
			if !UpdateByIdField(DBCAccounts, user.ID, "firstname", user.FirstName) {
				//update firstname failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.LastName) != 0 && user.LastName != userStored.LastName {
			if !UpdateByIdField(DBCAccounts, user.ID, "lastname", user.LastName) {
				//update lastname failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.Password) != 0 && user.Password != userStored.Password {
			if !UpdateByIdField(DBCAccounts, user.ID, "password", user.Password) {
				//update password failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.Birthday) != 0 && user.Birthday != userStored.Birthday {
			if !UpdateByIdField(DBCAccounts, user.ID, "birthday", user.Birthday) {
				//update birthday failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.Gender) != 0 && user.Gender != userStored.Gender {
			if !UpdateByIdField(DBCAccounts, user.ID, "gender", user.Gender) {
				//update gender failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if len(user.Comments) != 0 && user.Comments != userStored.Comments {
			if !UpdateByIdField(DBCAccounts, user.ID, "comments", user.Comments) {
				//update comments failed
				rsp.Code = RspBadRequest
				rsp.Reason = ReasonOperationFailed
				RespondJ(c, RspBadRequest, rsp)
				return NewError(ReasonOperationFailed)
			}
			changed = true
		}
		if user.CarePoints != nil {
			c.Logger().Debug("method: ", method, ",  CarePoints: ", user.CarePoints, " len: ", len(user.CarePoints))
			if method == "push" {
				if !PushNewCarePoint(DBCAccounts, user.ID, user.CarePoints) {
					rsp.Code = RspBadRequest
					rsp.Reason = ReasonOperationFailed
					RespondJ(c, RspBadRequest, rsp)
					return NewError(ReasonOperationFailed)
				}
				changed = true
			} else if method == "pull" {
				if !PullNewCarePoint(DBCAccounts, user.ID, user.CarePoints) {
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
	return nil
}

func HandleGetUsers(c echo.Context, user *User, debug bool) error {
	rsp := &Response{RspOK, ReasonSuccess, nil, 0}
	if debug {
		//return all users
		users, err := FindAllUsers()
		if err == nil {
			c.Logger().Debug("Found users size:", len(users))
			rsp.Data = users
			rsp.Count = len(users)
			RespondJ(c, RspOK, rsp)
			return nil
		}
		c.Logger().Debug("FindAlluser failed!!err:", err.Error())
	} else if userFound := FindUserWithID(user.ID); userFound != nil {
		c.Logger().Debug("Found users :", JsonToString(userFound))
		rsp.Data = userFound
		rsp.Count = 1
		RespondJ(c, RspOK, rsp)
		return nil
	}
	rsp.Code = RspBadRequest
	rsp.Reason = ReasonNotFound
	RespondJ(c, RspBadRequest, rsp)
	return NewError("Not found any")
}

func LoginValidate(c echo.Context, username, password string) (*User, error) {
	m := make(map[string]string)
	m["key1"] = "associatedId"
	m["value1"] = username
	m["key2"] = "phonenumber"
	m["value2"] = username
	m["key3"] = "email"
	m["value3"] = username
	m["key4"] = "password"
	m["value4"] = password
	c.Logger().Debug("ValidateAccount :", m)
	if users, err := FindUsersWithPW(DBCAccounts, m); users != nil && err == nil {
		//found user registered
		c.Logger().Debug("Found users:", JsonToString(users))
		//how about users is more than one element
		return users[0], nil
	}
	return nil, NewError(ReasonNotFound)
}

func GetAccountIDViaUserID(c echo.Context, userid string) (*User, error) {
	m := make(map[string]string)
	if len(userid) != 0 {
		m["key1"] = "associatedId"
		m["value1"] = userid
	}
	if len(userid) != 0 {
		m["key2"] = "phonenumber"
		m["value2"] = userid
	}
	if len(userid) != 0 {
		m["key3"] = "email"
		m["value3"] = userid
	}
	c.Logger().Debug("map side:", len(m))
	if usersFound, err := FindUsersWithFeild(DBCAccounts, m); err == nil {
		if len(usersFound) != 0 {
			c.Logger().Debug("Found users size: ", len(usersFound))
			//how about users is more than one element
			return usersFound[0], err
		}
	}
	return nil, NewError("Not found")
}
