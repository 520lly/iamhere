package services

import (
	"github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MsgService struct {
}

func (this *MsgService) handleMessagesPost(msg *Info.Message) {

	if err := decodeBody(r, &m); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read msg from request!! error:", err)
		return
	}
	if len(m.UserID) == 0 {
		responseHandleMessage(w, r, http.StatusBadRequest, "UserID is empty", nil)
		return
	} else if len(m.Content) == 0 {
		responseHandleMessage(w, r, http.StatusBadRequest, "Content is empty", nil)
		return
	} else if m.Latitude >= LatitudeMaximum || m.Latitude <= LatitudeMinimum {
		responseHandleMessage(w, r, http.StatusBadRequest, "latitude is out of range", nil)
		return
	} else if m.Longitude >= LongitudeMaximum || m.Longitude <= LongitudeMinimum {
		responseHandleMessage(w, r, http.StatusBadRequest, "longitude is out of range", nil)
		return
	} else if m.ExpiryTime == 0 {
		m.ExpiryTime = time.Now().Unix()
	}
	log.Println("msg longitude: ", m.Longitude, "latitude: ", m.Latitude, "altitude: ", m.Altitude)

	m.TimeStamp = time.Now()
	m.Location.Coordinates = []float64{m.Longitude, m.Latitude}
	m.Location.Type = "Point"
	m.ID = bson.NewObjectId()
	err := c.Insert(m)
	if err != nil {
		responseHandleMessage(w, r, http.StatusInternalServerError, ReasonInsertFailure, nil)
		return
	}
	// ensure
	// Creating the indexes
	index := mgo.Index{
		Key: []string{"$2dsphere:location"},
	}
	err = c.EnsureIndex(index)
	if err != nil {
		log.Println("There is index error")
		respondErr(w, r, http.StatusBadRequest, err, nil)
	}
	w.Header().Set("Location", "message/"+m.ID.Hex())
	responseHandleMessage(w, r, RspOK, ReasonSuccess, nil)
}
