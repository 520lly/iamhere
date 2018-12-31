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
	MessageTypeOcean string = "Ocean"
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
	Author      string        `json:"author" bson:"author"`
	UserDefAddr string        `json:"userdefaddr" bson:"userdefaddr"`
	ExpiryTime  int64         `json:"expirytime"`
	Altitude    float64       `json:"altitude"`
	Location    GeoJson       `bson:"location" json:"location"`
	APIKey      string        `json:"apikey"`
	Latitude    float64       `json:"latitude"`
	Longitude   float64       `json:"longitude"`
	TimeStamp   int64         `json:"timestamp"`
	LikeCount   int32         `json:"likecount"`
	Recommend   bool          `json:"recommend"`
	Color       int32         `json:color`
	Available   bool          `json:available`
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
	case "PUT":
		log.Println("PUT")
		s.handleMessagesPut(w, r)
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
	var msgs []*Message
	if p.HasID() {
		// get specific messages
		log.Println("ID ", p.ID)
		var err error = nil
		if err = c.FindId(bson.ObjectIdHex(p.ID)).All(&msgs); err != nil {
			log.Println("error ", string(err.Error()))
			responseHandleMessage(w, r, RspOK, err.Error(), nil)
			return
		}
		responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
		return
	} else {
		// get all messages
		q = c.Find(nil)
	}
	//get all list for debugging
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
	areaid := r.URL.Query().Get("areaid")
	log.Println("areaid=", areaid)
	if len(areaid) != 0 {
		if areaid == "Ocean" {
			//return some of messages from Ocean
			err := session.DB("iamhere").C("msgcoean").Find(nil).All(&msgs)
			if err == nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
			log.Println("msgcoean msgs size", len(msgs))
		} else {
			err := c.Find(bson.M{"areaid": areaid}).All(&msgs)
			log.Println("msgs size", len(msgs))
			if err == nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
		}
	}
	msgLon := r.URL.Query().Get("longitude")
	msgLat := r.URL.Query().Get("latitude")
	msgAlt := r.URL.Query().Get("altitude")
	msgRad := r.URL.Query().Get("radius")
	log.Println("URL Quary Message longitude: ", msgLon, "latitude: ", msgLat, "altitude: ", msgAlt, "radius: ", msgRad)
	if len(msgLat) != 0 || len(msgLon) != 0 || len(msgAlt) != 0 || len(msgRad) != 0 {
		if !(len(msgLat) != 0 && len(msgLon) != 0) {
			responseHandleMessage(w, r, http.StatusBadRequest, "latitude,longitude must be valid at the same time", nil)
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
		log.Println("msg longitude: ", alon, "latitude: ", alat)
		//find
		area := s.findAreaWithLocation(alon, alat)
		if area != nil {
			//find
			err := c.Find(bson.M{
				"location": bson.M{
					"$nearSphere": bson.M{
						"$geometry": bson.M{
							"Type":        "Point",
							"coordinates": []float64{alon, alat},
						},
						"$maxDistance": area.Radius,
					},
				},
			}).All(&msgs)
			if err != nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
			log.Println("msgs size", len(msgs))
		} else {
			//return some of messages from Ocean
			err := session.DB("iamhere").C("msgcoean").Find(nil).All(&msgs)
			if err != nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
			log.Println("msgcoean msgs size", len(msgs))
		}
	} else {
		var m Message
		if err := decodeBody(r, &m); err != nil {
			respondErr(w, r, http.StatusBadRequest, "failed to read msg from request!! error:", err)
			return
		}
		log.Println("Request Quary data longitude: ", m.Longitude, "latitude: ", m.Latitude)
		if !checkInRangefloat64(m.Longitude, LongitudeMinimum, LongitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "longitude is out of range", nil)
			return
		}
		if !checkInRangefloat64(m.Latitude, LatitudeMinimum, LatitudeMaximum) {
			responseHandleMessage(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		area := s.findAreaWithLocation(m.Longitude, m.Latitude)
		if area != nil {
			//find
			err := c.Find(bson.M{
				"location": bson.M{
					"$nearSphere": bson.M{
						"$geometry": bson.M{
							"Type":        "Point",
							"coordinates": []float64{m.Longitude, m.Latitude},
						},
						"$maxDistance": area.Radius,
					},
				},
			}).All(&msgs)
			if err != nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
			log.Println("msgs size", len(msgs))
		} else {
			//return some of messages from Ocean
			err := session.DB("iamhere").C("msgcoean").Find(nil).All(&msgs)
			if err != nil {
				responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
				return
			}
			log.Println("msgcoean msgs size", len(msgs))
		}
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

func (s *Server) handleMessagesPut(w http.ResponseWriter, r *http.Request) {
	log.Println("handleMessagesPut")
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("messages")
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusBadRequest, "failed to update without msg id", nil)
		return
	}
	log.Println("ID=", p.ID)
	likecount := r.URL.Query().Get("likecount")
	log.Println("likecount=", likecount)
	if len(likecount) != 0 {
		lc, err := strconv.ParseInt(likecount, 0, 32)
		if err != nil {
			responseHandleMessage(w, r, RspFailed, ReasonOperationFailed, nil)
			return
		}
		log.Println("likecount=", lc)
		colQuerier := bson.M{"_id": bson.ObjectIdHex(p.ID)}
		change := bson.M{"$set": bson.M{"likecount": lc}}
		err = c.Update(colQuerier, change)
		if err != nil {
			log.Println("error=", err.Error())
			responseHandleMessage(w, r, RspFailed, ReasonOperationFailed, nil)
			return
		}
	}
	recommend := r.URL.Query().Get("recommend")
	log.Println("recommend=", recommend)
	if len(recommend) != 0 {
		rc, err := strconv.ParseBool(recommend)
		if err != nil {
			responseHandleMessage(w, r, RspFailed, ReasonOperationFailed, nil)
			return
		}
		log.Println("recommend=", rc)
		colQuerier := bson.M{"_id": bson.ObjectIdHex(p.ID)}
		change := bson.M{"$set": bson.M{"recommend": rc}}
		err = c.Update(colQuerier, change)
		if err != nil {
			responseHandleMessage(w, r, RspFailed, ReasonOperationFailed, nil)
			return
		}
	}
	var msgs []*Message
	if err := c.FindId(bson.ObjectIdHex(p.ID)).All(&msgs); err != nil {
		log.Println("error ", string(err.Error()))
		responseHandleMessage(w, r, RspOK, err.Error(), nil)
		return
	}
	responseHandleMessage(w, r, RspOK, ReasonSuccess, &msgs)
}

func (s *Server) handleMessagesPost(w http.ResponseWriter, r *http.Request) {
	log.Println("handleMessagesPost")
	var m Message
	var geoEnable bool = false
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
	} else if m.ExpiryTime == 0 {
		m.ExpiryTime = time.Now().Unix()
	}
	var area *Area = nil
	if m.Latitude >= LatitudeMaximum || m.Latitude <= LatitudeMinimum {
		responseHandleMessage(w, r, http.StatusBadRequest, "latitude is out of range", nil)
		return
	} else if m.Longitude >= LongitudeMaximum || m.Longitude <= LongitudeMinimum {
		responseHandleMessage(w, r, http.StatusBadRequest, "longitude is out of range", nil)
		return
	} else {
		log.Println("msg longitude: ", m.Longitude, "latitude: ", m.Latitude, "altitude: ", m.Altitude)
		m.Location.Coordinates = []float64{m.Longitude, m.Latitude}
		m.Location.Type = "Point"
		if len(m.AreaID) != 0 {
			//use areaid
			log.Println("m.AreaID=", m.AreaID)
			if m.AreaID != "Ocean" {
				area = s.findAreaWithID(m.AreaID)
			} else if m.Latitude == 0 && m.Latitude == 0 {
				responseHandleMessage(w, r, http.StatusBadRequest, "AreaID or (longitude and latitude) must be available", nil)
				return
			}
		} else {
			if m.Latitude == 0 && m.Latitude == 0 {
				responseHandleMessage(w, r, http.StatusBadRequest, "AreaID or (longitude and latitude) must be available", nil)
				return
			} else {
				area = s.findAreaWithLocation(m.Longitude, m.Latitude)
				if area != nil {
					m.AreaID = string(bson.ObjectId(area.ID).Hex())
					geoEnable = true
				} else {
					//Ocean message
					m.AreaID = "Ocean"
					log.Println("m.AreaID=", m.AreaID)
				}
			}
		}
	}
	m.TimeStamp = time.Now().Unix()
	m.ID = bson.NewObjectId()
	session := s.db.Copy()
	defer session.Close()
	var c *mgo.Collection
	if area != nil {
		log.Println("area.ID=", area.ID)
		c = session.DB("iamhere").C("messages")
		err := c.Insert(m)
		if err != nil {
			responseHandleMessage(w, r, http.StatusInternalServerError, ReasonInsertFailure, nil)
			return
		}
	} else {
		if len(m.AreaID) != 0 {
			//it's a none area belonging message
			m.AreaID = MessageTypeOcean
			log.Println("area.ID=", MessageTypeOcean)
			c = session.DB("iamhere").C("msgcoean")
			err := c.Insert(m)
			if err != nil {
				responseHandleMessage(w, r, http.StatusInternalServerError, ReasonInsertFailure, nil)
				return
			}
		}
	}
	log.Println("geoEnable=", geoEnable)
	if geoEnable {
		// ensure
		// Creating the indexes
		index := mgo.Index{
			Key: []string{"$2dsphere:location"},
		}
		err := c.EnsureIndex(index)
		if err != nil {
			log.Println("There is index error")
			respondErr(w, r, http.StatusBadRequest, err, nil)
		}
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
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete message without Message ID.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		c := session.DB("iamhere").C("msgcoean")
		if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
			respondErr(w, r, http.StatusInternalServerError, "failed to delete message. error:", err)
			return
		}
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
