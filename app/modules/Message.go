package modules

import (
	"github.com/520lly/iamhere/app/modules"
	"time"
)

type Message struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	AreaID      string        `json:"areaid" bson:"areaid"`
	UserID      string        `json:"userid" bson:"userid"`
	Content     string        `json:"content" bson:"content"`
	UserDefAddr string        `json:"userdefaddr" bson:"userdefaddr"`
	ExpiryTime  int64         `json:"expirytime"`
	Altitude    float64       `json:"altitude"`
	Location    Info.GeoJson  `bson:"location" json:"location"`
	APIKey      string        `json:"apikey"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	TimeStamp   time.Time
}
