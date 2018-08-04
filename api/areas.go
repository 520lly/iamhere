package main

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type area struct {
	ID      bson.ObjectId  `bson:"_id" json:"id"`
	Title   string         `json:"title" bson:"title"`
	Options []string       `json:"options"`
	Results map[string]int `json:"results,omitempty"`
	APIKey  string         `json:"apikey"`
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
	c := session.DB("ballots").C("areas")
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific area
		q = c.FindId(bson.ObjectIdHex(p.ID))
	} else {
		// get all areas
		q = c.Find(nil)
	}
	var result []*area
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, &result)
}

func (s *Server) handleareasPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("areas")
	var p area
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read area from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		p.APIKey = apikey
	}
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert area", err)
		return
	}
	w.Header().Set("Location", "areas/"+p.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

func (s *Server) handleareasDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("ballots").C("areas")
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all areas.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(p.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete area", err)
		return
	}
	respond(w, r, http.StatusOK, nil) // ok
}
