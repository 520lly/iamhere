package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	RspOK     int = 0
	RspFailed int = -1
)
const (
	ReasonSuccess         string = "Success"
	ReasonFailureParam    string = "Wrong parameter"
	ReasonMissingParam    string = "Missing parameter"
	ReasonFailureAPIKey   string = "Wrong APIKey"
	ReasonFailueGeneral   string = "Failure in general"
	ReasonDuplicate       string = "Parameter duplicated"
	ReasonInsertFailure   string = "Insert failed"
	ReasonWrongPw         string = "Wrong Password "
	ReasonOperationFailed string = "Operation Failure "
)

type Message struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	AreaID      string        `json:"areaid" bson:"areaid"`
	UserID      string        `json:"userid" bson:"userid"`
	Content     string        `json:"content" bson:"content"`
	UserDefAddr string        `json:"userdefaddr" bson:"userdefaddr"`
	ExpiryTime  int64         `json:"expirytime"`
	Altitude    float64       `json:"altitude"`
	Location    GeoJson       `bson:"location" json:"location"`
	APIKey      string        `json:"apikey"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	TimeStamp   time.Time
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleMessagesGet(w, r)
		return
	case "POST":
		log.Println("POST")
		s.handleMessagesPost(w, r)
		return
	case "DELETE":
		s.handleMessagesDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleMessagesGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("messages")
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific messages
		log.Println("ID ", p.HasID())
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all messages
		q = c.Find(nil)
	}
	//get all list for debugging
	var msgs []*Message
	debug := r.URL.Query().Get("debug")
	log.Println("debug=", debug)
	if len(debug) != 0 {
		if err := q.All(&msgs); err != nil {
			respondErr(w, r, http.StatusInternalServerError, err)
			return
		}
		log.Println("msgs size", len(msgs))
		responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
		return
	}
	msgLon := r.URL.Query().Get("longitude")
	msgLat := r.URL.Query().Get("latitude")
	msgAlt := r.URL.Query().Get("altitude")
	msgRad := r.URL.Query().Get("radius")
	log.Println("Message longitude: ", msgLon, "latitude: ", msgLat, "altitude: ", msgAlt, "radius: ", msgRad)
	if len(msgLat) != 0 || len(msgLon) != 0 || len(msgAlt) != 0 || len(msgRad) != 0 {
		if !(len(msgLat) != 0 && len(msgLon) != 0 && len(msgAlt) != 0 && len(msgRad) != 0) {
			responseHandleMessage(w, r, http.StatusBadRequest, "latitude,longitude, altitude, radius must be valid at the same time", nil)
			return
		}
		alon, err := strconv.ParseFloat(msgLon, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alon, LongitudeMinimum, LongitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "longitude is out of range", nil)
			return
		}
		alat, err := strconv.ParseFloat(msgLat, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alat, LatitudeMinimum, LatitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		aalt, err := strconv.ParseFloat(msgAlt, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(aalt, AltitudeMinium, AltitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "altitude is out of range", nil)
			return
		}
		arad, err := strconv.ParseFloat(msgRad, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alat, LatitudeMinimum, LatitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		log.Println("msg longitude: ", alon, "latitude: ", alat, "altitude: ", aalt, "radius: ", arad)
		//find
		err = c.Find(bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"Type":        "Point",
						"coordinates": []float64{alon, alat},
					},
					"$maxDistance": arad,
				},
			},
		}).All(&msgs)
	} else {
		/*check area id??*/
	}
	//check ExpiryTime and remove none expiry messages
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].ExpiryTime > time.Now().Unix() {
			msgs = append(msgs[:i], msgs[i+1:]...)
		}
	}
	log.Println("msgs size", len(msgs))
	responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
}

func (s *Server) handleMessagesPost(w http.ResponseWriter, r *http.Request) {
	log.Println("handleMessagesPost")
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("messages")
	var m Message
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

func (s *Server) handleMessagesDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	// Collection Message
	c := session.DB("iamhere").C("messages")
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all messages.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete message", err)
		return
	}
	responseHandleMessage(w, r, RspOK, ReasonSuccess, nil)
}

func responseHandleMessage(w http.ResponseWriter, r *http.Request, code int, reason string, msgs *[]*Message) {
	type response struct {
		Code   int         `json:"code"`
		Reason string      `json:"reason"`
		Data   *[]*Message `json:"data"`
		Count  int         `json:"count"`
	}
	result := &response{
		Code:   code,
		Reason: reason,
		Data:   msgs,
		Count:  0}
	if msgs != nil {
		result.Count = len(*msgs)
	}
	respond(w, r, http.StatusOK, &result)
}
