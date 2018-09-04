package modules

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	AssociatedId string        `json:"associatedId" bson:"associatedId"`
	PhoneNumber  string        `json:"phonenumber" bson:"phonenumber"`
	Email        string        `json:"email" bson:"email"`
	NickName     string        `json:"nickname,omitempty"`
	FirstName    string        `json:"firstname,omitempty"`
	LastName     string        `json:"lastname,omitempty"`
	Password     string        `json:"password"`
	Birthday     string        `json:"birthday,omitempty"`
	Gender       string        `json:"gender,omitempty"`
	Comments     string        `json:"comments,omitempty"`
	TimeStamp    int64         `json:"timestamp"`
}
