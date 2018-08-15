package main

import (
	//"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LatitudeMinimum float32 = -90 //Latitude Minimum
	LatitudeMaximum float32 = 90  //Latitude Maximum
)

const (
	LongitudeMinimum float32 = -180 //Longitude Minimum
	LongitudeMaximum float32 = 180  //Longitude Maximum
)

const (
	AltitudeMinium  float32 = -10000 //Altitude Minimum
	AltitudeMaximum float32 = 10000  //Altitude Maximum
)

const (
	RadiusMinimum float32 = 10    // Radius Minimum meters
	RadiusMaximum float32 = 10000 // Radius Maximum
)

const (
	CategoryMinimum int = -1               //Category Minimum
	CategorySystem  int = 0                //defined by system
	CategoryUser    int = 1                //defined by user
	CategoryMaximum int = CategoryUser + 1 //Category Maximum
)

const (
	TypeMinimum    int = -1                 //Area type Minimum
	TypeNatural    int = 0                  //Area type is Natural
	TypeHistorical int = 1                  //Area type is Historical
	TypeHumanities int = 2                  //Area type is Humanities
	TypeMaximum    int = TypeHumanities + 1 //Area type Maximum
)

type Area struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description"`
	Address1    string        `json:"address1"`
	Address2    string        `json:"address2"`
	Category    int           `json:"category"`
	Type        int           `json:"type"`
	Latitude    float32       `json:"latitude"`
	Longitude   float32       `json:"longitude"`
	Altitude    float32       `json:"altitude"`
	Radius      float32       `json:"radius"` //meter
	APIKey      string        `json:"apikey"`
}

func (s *Server) handleareas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleareasGet(w, r)
		return
	case "POST":
		s.handleareasPost(w, r)
		return
	case "DELETE":
		s.handleareasDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleareasGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("areas")
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific area
		log.Println("ID ", p.HasID())
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all areas
		q = c.Find(nil)
	}
	var areas []*Area
	if err := q.All(&areas); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	areaType := r.URL.Query().Get("type")
	if len(areaType) != 0 {
		t, err := strconv.Atoi(areaType)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		log.Println("area Type is ", t)
		if !checkInRange(t, TypeMinimum, TypeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "type is out of range", nil)
			return
		}
		for i, area := range areas {
			if area.Type != t {
				areas = append(areas[:i], areas[i+1:]...)
			}
		}
	}
	areaCategory := r.URL.Query().Get("category")
	if len(areaCategory) != 0 {
		ac, err := strconv.Atoi(areaCategory)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		log.Println("area category is ", ac)
		if !checkInRange(ac, CategoryMinimum, CategoryMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "category is out of range", nil)
			return
		}
		for i, area := range areas {
			if area.Category != ac {
				areas = append(areas[:i], areas[i+1:]...)
			}
		}
	}
	areaLat := r.URL.Query().Get("latitude")
	areaLon := r.URL.Query().Get("longitude")
	areaAlt := r.URL.Query().Get("altitude")
	areaRad := r.URL.Query().Get("radius")
	if len(areaLat) != 0 || len(areaLon) != 0 || len(areaAlt) != 0 || len(areaRad) != 0 {
		if !(len(areaLat) != 0 && len(areaLon) != 0 && len(areaAlt) != 0 && len(areaRad) != 0) {
			responseHandleAreas(w, r, http.StatusBadRequest, "latitude,longitude, altitude, radius must be valid at the same time", nil)
			return
		}
		alat, err := strconv.ParseFloat(areaLat, 32)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		alat32 := float32(alat)
		alon, err := strconv.ParseFloat(areaLon, 32)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		alon32 := float32(alon)
		aalt, err := strconv.ParseFloat(areaAlt, 32)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		aalt32 := float32(aalt)
		arad, err := strconv.ParseFloat(areaRad, 32)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		arad32 := float32(arad)

		log.Println("area latitude: ", alat32, "longitude: ", alon32, "altitude: ", aalt32, "radius: ", arad32)
		if !checkInRangeFloat32(alat32, LatitudeMinimum, LatitudeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		for i, area := range areas {
			if area.Latitude != alat32 {
				areas = append(areas[:i], areas[i+1:]...)
			}
		}
	}
	log.Println("areas size", len(areas))
	responseHandleAreas(w, r, RspOK, ReasonSuccess, &areas)
}

func (s *Server) handleareasPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("areas")
	var p Area
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read area from request", err)
		return
	}
	if len(p.Description) == 0 {
		responseHandleAreas(w, r, http.StatusBadRequest, "description is empty", nil)
		return
	} else if len(p.Address1) == 0 {
		responseHandleAreas(w, r, http.StatusBadRequest, "address1 is empty", nil)
		return
	} else if p.Category >= CategoryMaximum || p.Category <= CategoryMinimum {
		responseHandleAreas(w, r, http.StatusBadRequest, "category is out of range", nil)
		return
	} else if p.Type >= TypeMaximum || p.Type <= TypeMinimum {
		responseHandleAreas(w, r, http.StatusBadRequest, "type is out of range", nil)
		return
	} else if p.Latitude >= LatitudeMaximum || p.Latitude <= LatitudeMinimum {
		responseHandleAreas(w, r, http.StatusBadRequest, "latitude is out of range", nil)
		return
	} else if p.Longitude >= LongitudeMaximum || p.Longitude <= LongitudeMinimum {
		responseHandleAreas(w, r, http.StatusBadRequest, "longitude is out of range", nil)
		return
	} else if p.Radius >= RadiusMaximum || p.Radius <= RadiusMinimum {
		responseHandleAreas(w, r, http.StatusBadRequest, "radius is out of range", nil)
		return
	}
	apikey, ok := APIKey(r.Context())
	if !ok {
		responseHandleAreas(w, r, RspFailed, ReasonFailureAPIKey, nil)
		return
	}
	p.APIKey = apikey
	var q *mgo.Query
	// get all areas
	q = c.Find(nil)
	var areas []*Area
	if err := q.All(&areas); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	log.Println("area.name is %s", string(p.Name))
	for _, area := range areas {
		if area.Name == p.Name {
			responseHandleAreas(w, r, http.StatusInternalServerError, ReasonDuplicate, nil)
			return
		}
	}
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		responseHandleAreas(w, r, http.StatusInternalServerError, ReasonInsertFailure, nil)
		return
	}
	w.Header().Set("Location", "areas/"+p.ID.Hex())
	responseHandleAreas(w, r, RspOK, ReasonSuccess, nil)
}

func (s *Server) handleareasDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("areas")
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all areas.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete area", err)
		return
	}
	responseHandleAreas(w, r, RspOK, ReasonSuccess, nil)
}

func responseHandleAreas(w http.ResponseWriter, r *http.Request, code int, reason string, areas *[]*Area) {
	type response struct {
		Code   int      `json:"code"`
		Reason string   `json:"reasone"`
		Data   *[]*Area `json:"data"`
	}
	result := &response{
		Code:   code,
		Reason: reason,
		Data:   areas}
	respond(w, r, http.StatusOK, &result)
}

func checkInRange(num int, bottom int, top int) (ret bool) {
	return num > bottom && num < top
}

func checkInRangeFloat32(num float32, bottom float32, top float32) (ret bool) {
	return num > bottom && num < top
}
