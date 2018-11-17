package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	LatitudeMinimum float64 = -90 //Latitude Minimum
	LatitudeMaximum float64 = 90  //Latitude Maximum
)

const (
	LongitudeMinimum float64 = -180 //Longitude Minimum
	LongitudeMaximum float64 = 180  //Longitude Maximum
)

const (
	AltitudeMinium  float64 = -10000 //Altitude Minimum
	AltitudeMaximum float64 = 10000  //Altitude Maximum
)

const (
	RadiusMinimum float64 = 10    // Radius Minimum meters
	RadiusMaximum float64 = 10000 // Radius Maximum
)

const (
	CategoryMinimum int = -1               //Category Minimum
	CategoryFixed   int = 1                //defined by system for fixed areas
	CategoryOcean   int = 2                //defined by system for ocean
	CategoryIsland  int = 4                //defined by system for island
	CategoryCloud   int = 8                //defined by system for cloud
	CategoryUser    int = 16               //defined by user
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
	Altitude    float64       `json:"altitude"`
	Radius      float64       `json:"radius"` //meter
	Location    GeoJson       `bson:"location" json:"location"`
	APIKey      string        `json:"apikey"`
	Longitude   float64       `json:"longitude"`
	Latitude    float64       `json:"latitude"`
	TimeStamp   int64         `json:"timestamp"`
}

type GeoJson struct {
	Type        string    `josn:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func (s *Server) handleAreas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleAreasGet(w, r)
		return
	case "POST":
		s.handleAreasPost(w, r)
		return
	case "DELETE":
		s.handleAreasDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleAreasGet(w http.ResponseWriter, r *http.Request) {
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
	//get all list for debugging
	var areas []*Area
	debug := r.URL.Query().Get("debug")
	log.Println("debug=", debug)
	if len(debug) != 0 {
		if err := q.All(&areas); err != nil {
			respondErr(w, r, http.StatusInternalServerError, err)
			return
		}
		responseHandleAreas(w, r, RspOK, ReasonSuccess, &areas)
		return
	}
	areaLon := r.URL.Query().Get("longitude")
	areaLat := r.URL.Query().Get("latitude")
	areaAlt := r.URL.Query().Get("altitude")
	areaRad := r.URL.Query().Get("radius")
	log.Println("area longitude: ", areaLon, "latitude: ", areaLat, "altitude: ", areaAlt, "radius: ", areaRad)
	if len(areaLat) != 0 || len(areaLon) != 0 || len(areaAlt) != 0 || len(areaRad) != 0 {
		if !(len(areaLat) != 0 && len(areaLon) != 0 && len(areaAlt) != 0 && len(areaRad) != 0) {
			responseHandleAreas(w, r, http.StatusBadRequest, "latitude,longitude, altitude, radius must be valid at the same time", nil)
			return
		}
		alon, err := strconv.ParseFloat(areaLon, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alon, LongitudeMinimum, LongitudeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "longitude is out of range", nil)
			return
		}
		alat, err := strconv.ParseFloat(areaLat, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alat, LatitudeMinimum, LatitudeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		aalt, err := strconv.ParseFloat(areaAlt, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(aalt, AltitudeMinium, AltitudeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "altitude is out of range", nil)
			return
		}
		arad, err := strconv.ParseFloat(areaRad, 64)
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkInRangefloat64(alat, LatitudeMinimum, LatitudeMaximum) {
			responseHandleAreas(w, r, http.StatusBadRequest, "latitude is out of range", nil)
			return
		}
		log.Println("area longitude: ", alon, "latitude: ", alat, "altitude: ", aalt, "radius: ", arad)
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
		}).All(&areas)
	} else {
		if err := q.All(&areas); err != nil {
			responseHandleAreas(w, r, http.StatusNotFound, ReasonSuccess, nil)
			return
		}
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
		for i := len(areas) - 1; i >= 0; i-- {
			if areas[i].Type != t {
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
		for i := len(areas) - 1; i >= 0; i-- {
			log.Println("areas[i].Category is ", areas[i].Category)
			if ac&1 == 0 { //request category excluded CategoryFixed
				if areas[i].Category == CategoryFixed {
					areas = append(areas[:i], areas[i+1:]...)
					if len(areas) <= 0 {
						break
					}
				}
			}
			if ac&(1<<1) == 0 { //request category excluded CategoryOcean
				if areas[i].Category == CategoryOcean {
					areas = append(areas[:i], areas[i+1:]...)
					if len(areas) <= 0 {
						break
					}
				}
			}
			if ac&(1<<2) == 0 { //request category excluded CategoryIsland
				if areas[i].Category == CategoryIsland {
					areas = append(areas[:i], areas[i+1:]...)
					if len(areas) <= 0 {
						break
					}
				}
			}
			if ac&(1<<3) == 0 { //request category excluded CategoryCloud
				if areas[i].Category == CategoryCloud {
					areas = append(areas[:i], areas[i+1:]...)
					if len(areas) <= 0 {
						break
					}
				}
			}
			if ac&(1<<4) == 0 { //request category excluded CategoryUser
				if areas[i].Category == CategoryUser {
					areas = append(areas[:i], areas[i+1:]...)
					if len(areas) <= 0 {
						break
					}
				}
			}
		}
	}
	log.Println("areas size", len(areas))
	responseHandleAreas(w, r, RspOK, ReasonSuccess, &areas)
}

func (s *Server) handleAreasPost(w http.ResponseWriter, r *http.Request) {
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
	}
	//else if p.Radius >= RadiusMaximum || p.Radius <= RadiusMinimum {
	//responseHandleAreas(w, r, http.StatusBadRequest, "radius is out of range", nil)
	//return
	//}
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
	p.Location.Coordinates = []float64{p.Longitude, p.Latitude}
	p.Location.Type = "Point"
	p.ID = bson.NewObjectId()
	p.TimeStamp = time.Now().Unix()
	err := c.Insert(p)
	if err != nil {
		responseHandleAreas(w, r, http.StatusInternalServerError, ReasonInsertFailure, nil)
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
		//panic(err)
	}
	w.Header().Set("Location", "areas/"+p.ID.Hex())
	responseHandleAreas(w, r, RspOK, ReasonSuccess, nil)
}

func (s *Server) handleAreasDelete(w http.ResponseWriter, r *http.Request) {
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

func checkInRange(num int, bottom int, top int) (ret bool) {
	return num > bottom && num < top
}

func checkInRangefloat64(num float64, bottom float64, top float64) (ret bool) {
	return num > bottom && num < top
}

func (s *Server) findAreaWithLocation(long float64, lat float64) (ret *Area) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("areas")
	areas := findAllArea(session)
	if areas != nil {
		for _, area := range areas {
			var areaMatchs []*Area
			err := c.Find(bson.M{
				"location": bson.M{
					"$nearSphere": bson.M{
						"$geometry": bson.M{
							"Type":        "Point",
							"coordinates": []float64{long, lat},
						},
						"$maxDistance": area.Radius,
					},
				},
			}).All(&areaMatchs)
			if err != nil {
				return nil
			}
			if len(areaMatchs) == 1 && areaMatchs[0].ID == area.ID {
				areaMatchJ, err := json.Marshal(areaMatchs[0])
				if err != nil {
					log.Fatalf("JSON marshaling failed: %s", err)
				}
				log.Println("found Area:", string(areaMatchJ))
				return areaMatchs[0]
			}
		}
	}
	return nil
}

func (s *Server) findAreaWithID(id string) (ret *Area) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("areas")
	var areaMatchs []*Area
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).All(&areaMatchs)
	if err != nil {
		return nil
	}
	if len(areaMatchs) > 0 {
		return areaMatchs[0]
	}
	return nil
}

func findAllArea(session *mgo.Session) (ret []*Area) {
	c := session.DB("iamhere").C("areas")
	var q *mgo.Query
	q = c.Find(nil)
	var areas []*Area
	if err := q.All(&areas); err != nil {
		return nil
	}
	return areas
}

func responseHandleAreas(w http.ResponseWriter, r *http.Request, code int, reason string, areas *[]*Area) {
	type response struct {
		Code   int      `json:"code"`
		Reason string   `json:"reason"`
		Data   *[]*Area `json:"data"`
		Count  int      `json:"count"`
	}
	result := &response{
		Code:   code,
		Reason: reason,
		Data:   areas,
		Count:  0}
	if areas != nil {
		result.Count = len(*areas)
	}
	respond(w, r, http.StatusOK, &result)
}
