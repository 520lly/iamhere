package modules

import (
	"gopkg.in/mgo.v2/bson"
)

type ColorMotion int

const (
	Black ColorMotion = iota
	Green
	Red
	Pink
)

type AvailableLocally bool

const (
	CanBeSeen  AvailableLocally = true
	CantBeSeen AvailableLocally = false
)

type Message struct {
	ID          bson.ObjectId    `bson:"_id" json:"id"`
	AreaID      string           `json:"areaid" bson:"areaid"`
	UserID      string           `json:"userid" bson:"userid"`
	Content     string           `json:"content" bson:"content"`
	UserDefAddr string           `json:"userdefaddr" bson:"userdefaddr"`
	Author      string           `json:"author"`
	ExpiryTime  int64            `json:"expirytime"`
	Altitude    float64          `json:"altitude"`
	Latitude    float64          `json:"latitude"`
	Longitude   float64          `json:"longitude"`
	TimeStamp   int64            `json:"timestamp"`
	Location    GeoJson          `bson:"location" json:"location"`
	LikeCount   string           `json:"likecount"`
	Recommend   string           `json:"recommend"`
	Color       ColorMotion      `json:color`
	Available   AvailableLocally `json:available`   //this message will be available when other user is in the area
	LimitAccess AvailableLocally `json:limitaccess` //this message will be limited to access when other user is in the area the defalue is false
}
