package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type User struct {
	ID          bson.ObjectId `json:"id "bson:"_id"`
	NickName    string        `json:"nickname,omitempty"`
	Email       string        `json:"email,omitempty"`
	FirstName   string        `json:"firstname,omitempty"`
	LastName    string        `json:"lastname,omitempty"`
	PhoneNumber string        `json:"phonenumber,omitempty"`
	Birthday    string        `json:"birthday,omitempty"`
	Gender      string        `json:"gender,omitempty"`
	Comments    string        `json:"comments,omitempty"`
	APIKey      string        `json:"apikey"`
}

func (s *Server) handleAccounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleAccountsGet(w, r)
		return
	case "POST":
		s.handleAccountsPost(w, r)
		return
	case "DELETE":
		s.handleAccountsDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handleAccountsGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("users")
	var q *mgo.Query
	u := NewPath(r.URL.Path)
	debug := r.URL.Query().Get("debug")
	log.Println("debug=", debug)
	if len(debug) != 0 {
		//get all list for debugging
		q = c.Find(nil)
	}
	if u.HasID() {
		// get specific poll
		q = c.FindId(bson.ObjectIdHex(u.ID))
	} else {
		responseHandleAccounts(w, r, RspOK, ReasonMissingParam, nil)
		return
	}
	var users []*User
	if err := q.All(&users); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	responseHandleAccounts(w, r, RspOK, ReasonSuccess, &users)
}

func (s *Server) handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("users")
	var u User
	if err := decodeBody(r, &u); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read poll from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		u.APIKey = apikey
	}
	u.ID = bson.NewObjectId()
	if err := c.Insert(u); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert poll", err)
		return
	}
	w.Header().Set("Location", "users/"+u.ID.Hex())
	respond(w, r, http.StatusCreated, nil)
}

func (s *Server) handleAccountsDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	c := session.DB("iamhere").C("users")
	u := NewPath(r.URL.Path)
	if !u.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all users.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(u.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete poll", err)
		return
	}
	respond(w, r, http.StatusOK, nil) // ok
}

func responseHandleAccounts(w http.ResponseWriter, r *http.Request, code int, reason string, users *[]*User) {
	type response struct {
		Code   int      `json:"code"`
		Reason string   `json:"reasone"`
		Data   *[]*User `json:"data"`
		Count  int      `json:"count"`
	}
	result := &response{
		Code:   code,
		Reason: reason,
		Data:   users,
		Count:  0}
	if users != nil {
		result.Count = len(*users)
	}
	respond(w, r, http.StatusOK, &result)
}
