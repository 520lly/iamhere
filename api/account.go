package main

import (
	//"crypto/md5"
	//"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID          bson.ObjectId `json:"id bson:"_id"`
	Password    string        `json:"passward"`
	AssocitesId string        `json:"associatedId,omitempty" bson:"associatedId"`
	NickName    string        `json:"nickname,omitempty"`
	Email       string        `json:"email,omitempty" bson:"email"`
	FirstName   string        `json:"firstname,omitempty"`
	LastName    string        `json:"lastname,omitempty"`
	PhoneNumber string        `json:"phonenumber" bson:"phonenumber"`
	Birthday    string        `json:"birthday,omitempty"`
	Gender      string        `json:"gender,omitempty"`
	Comments    string        `json:"comments,omitempty"`
	APIKey      string        `json:"apikey"`
	TimeStamp   time.Time
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
	} else if u.HasID() {
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
		respondErr(w, r, http.StatusBadRequest, "failed to read user information from request", err)
		return
	}
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific user
		log.Println("ID ", p.HasID(), " =", p.ID, "passward=", u.Password)
		//pw := r.URL.Query().Get("passward")
		if len(u.Password) != 0 {
			var q *mgo.Query
			q = c.FindId(bson.ObjectIdHex(p.ID))
			var users []*User
			if err := q.All(&users); err != nil {
				respondErr(w, r, http.StatusInternalServerError, err)
				return
			}
			//res, err := json.Marshal(users)
			//if err != nil {
			//    log.Fatalf("JSON marshaling failed: %s", err)
			//}
			//log.Println("found user:", string(res))
			log.Println("length of users:", len(users))
			//more than one user found it's a fatal error
			if len(users) == 1 {
				log.Println("Password stored:", string(users[0].Password))
				log.Println("Password stored:", u.Password)
				//mast check pw
				if users[0].Password == u.Password {
					//todo update user information
					responseHandleAccounts(w, r, RspOK, ReasonSuccess, &users)
					return
				} else {
					responseHandleAccounts(w, r, RspFailed, ReasonWrongPw, nil)
					return
				}
				log.Println("NickName stored:", string(users[0].NickName))
			} else {
				responseHandleAccounts(w, r, RspFailed, ReasonWrongPw, nil)
				return
			}
		} else {
			responseHandleAccounts(w, r, RspFailed, ReasonMissingParam, nil)
			return
		}
	}
	if len(u.Password) == 0 {
		responseHandleAccounts(w, r, http.StatusBadRequest, "Password is empty", nil)
		return
	} else if len(u.AssocitesId) == 0 && len(u.Email) == 0 && len(u.PhoneNumber) == 0 {
		responseHandleAccounts(w, r, http.StatusBadRequest, "AssocitesId/Email/PhoneNumber must be valid at least one", nil)
		return
	}
	u.TimeStamp = time.Now()
	err := c.Find(bson.M{"name": "Tom"}).One(&user)
	if err != nil {
		panic(err)
	}
	apikey, ok := APIKey(r.Context())
	if ok {
		u.APIKey = apikey
	}
	u.ID = bson.NewObjectId()
	//has := md5.Sum(u.)
	if err := c.Insert(u); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert user", err)
		return
	}
	w.Header().Set("Location", "users/"+u.ID.Hex())
	responseHandleAccounts(w, r, http.StatusCreated, ReasonSuccess, nil)
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
