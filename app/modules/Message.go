package modules

import (
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	AreaID      string        `json:"areaid" bson:"areaid"`
	UserID      string        `json:"userid" bson:"userid"`
	Content     string        `json:"content" bson:"content"`
	UserDefAddr string        `json:"userdefaddr" bson:"userdefaddr"`
	Author      string        `json:"author"`
	ExpiryTime  int64         `json:"expirytime"`
	Altitude    float64       `json:"altitude"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	TimeStamp   int64         `json:"timestamp"`
	Location    GeoJson       `bson:"location" json:"location"`
	LikeCount   string        `json:"likecount"`
	Recommend   string        `json:"recommend"`
}
