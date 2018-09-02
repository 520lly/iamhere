package modules

import (
	"gopkg.in/mgo.v2/bson"
)

type Area struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Province    string        `json:"province" bson:"province"`
	City        string        `json:"city" bson:"city"`
	District    string        `json:"district" bson:"district"`
	Discription string        `json:"discription"`
	Address1    string        `json:"address1"`
	Address2    string        `json:"address2"`
	Category    int           `json:"category"`
	Type        int           `json:"type"`
	Longitude   float64       `json:"longitude"`
	Latitude    float64       `json:"latitude"`
	Altitude    float64       `json:"altitude"`
	Radius      float64       `json:"radius"` //meter
	Location    GeoJson       `bson:"location" json:"location"`
	TimeStamp   int64         `json:"timestamp"`
}
