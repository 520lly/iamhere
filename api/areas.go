package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	CategorySystem int = 0 //defined by system
	CategoryUser   int = 1 //defined by user
)

type Area struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `json:"name" bson:"name"`
	Address1  string        `json:"address1"`
	Address2  string        `json:"address2"`
	Category  int           `json:"category"`
	Type      int           `json:"type"`
	Latitude  float32       `json:"latitude"`
	Longitude float32       `json:"longitude"`
	Altitude  float32       `json:"altitude"`
	Radius    float32       `json:"radius"` //meter
	APIKey    string        `json:"apikey"`
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
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all areas
		q = c.Find(nil)
	}
	var result []*Area
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	responseHandleAreas(w, r, RspOK, ReasonSuccess, &result)
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
	apikey, ok := APIKey(r.Context())
	if !ok {
		responseHandleAreas(w, r, RspFailed, ReasonFailureAPIKey, nil)
	}
	p.APIKey = apikey
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert area", err)
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
