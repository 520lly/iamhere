package modules

import (
	"gopkg.in/mgo.v2/bson"
)

type Area struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name" valid:"alphanum,printableascii"`
	Province    string        `json:"province" bson:"province" valid:"alphanum,printableascii"`
	City        string        `json:"city" bson:"city" valid:"alphanum,printableascii"`
	District    string        `json:"district" bson:"district" valid:"alphanum,printableascii"`
	Discription string        `json:"discription" valid:"alphanum,printableascii"`
	Address1    string        `json:"address1" valid:"alphanum,printableascii"`
	Address2    string        `json:"address2" valid:"alphanum,printableascii"`
	Category    int           `json:"category" valid:"alphanum,printableascii"`
	Type        int           `json:"type" valid:"alphanum,printableascii"`
	Longitude   float64       `json:"longitude" valid:"alphanum,printableascii"`
	Latitude    float64       `json:"latitude"`
	Altitude    float64       `json:"altitude"`
	Radius      float64       `json:"radius"` //meter
	Location    GeoJson       `bson:"location" json:"location"`
	TimeStamp   int64         `json:"timestamp"`
}
