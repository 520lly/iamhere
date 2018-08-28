package main

import (
	//"crypto/md5"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Password     string        `json:"password"`
	AssociatedId string        `json:"associatedId,omitempty" bson:"associatedId"`
	NickName     string        `json:"nickname,omitempty"`
	Email        string        `json:"email,omitempty" bson:"email"`
	FirstName    string        `json:"firstname,omitempty"`
	LastName     string        `json:"lastname,omitempty"`
	PhoneNumber  string        `json:"phonenumber" bson:"phonenumber"`
	Birthday     string        `json:"birthday,omitempty"`
	Gender       string        `json:"gender,omitempty"`
	Comments     string        `json:"comments,omitempty"`
	APIKey       string        `json:"apikey"`
	TimeStamp    time.Time
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
	// Collection User
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
		if !bson.IsObjectIdHex(u.ID) {
			responseHandleAccounts(w, r, RspFailed, ReasonFailureParam, nil)
			return
		}
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
	//user iterator for avoiding massive memory usage
	iter := c.Find(nil).Iter()
	for iter.Next(&users) {
		fmt.Println(users)
	}
	responseHandleAccounts(w, r, RspOK, ReasonSuccess, &users)
}

func (s *Server) handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	// Collection User
	c := session.DB("iamhere").C("users")
	var u User
	if err := decodeBody(r, &u); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read user information from request", err)
		return
	}
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific user
		log.Println("ID ", p.HasID(), " =", p.ID, "password=", u.Password)
		//if len(u.Password) != 0 {
		user := User{}
		if err := c.FindId(bson.ObjectIdHex(p.ID)).One(&user); err != nil {
			responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
			return
		}
		userj, err := json.Marshal(user)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		log.Println("found user:", string(userj))
		log.Println("Password stored:", string(user.Password), "Password received:", u.Password)
		//must check pw????
		//if user.Password == u.Password {
		//update user email
		if len(u.Email) != 0 && u.Email != user.Email {
			log.Println("Update user email! new email is:", u.Email)
			b, err := updateField(c, p.ID, "email", u.Email)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty email or same email:", u.Email)
		}
		//update user password
		if len(u.Password) != 0 && u.Password != user.Password {
			log.Println("Update user password! new password is:", u.Password)
			b, err := updateField(c, p.ID, "password", u.Password)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty password or same password")
		}
		//update user AssociatedId
		if len(u.AssociatedId) != 0 && u.AssociatedId != user.AssociatedId {
			log.Println("Update user AssociatedId! new AssociatedId is:", u.AssociatedId)
			b, err := updateField(c, p.ID, "associatedId", u.AssociatedId)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty AssociatedId or same AssociatedId")
		}
		//update user NickName
		if len(u.NickName) != 0 && u.NickName != user.NickName {
			log.Println("Update user NickName! new NickName is:", u.NickName)
			b, err := updateField(c, p.ID, "nickname", u.NickName)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty NickName or same NickName")
		}
		//update user FirstName
		if len(u.FirstName) != 0 && u.FirstName != user.FirstName {
			log.Println("Update user FirstName! new FirstName is:", u.FirstName)
			b, err := updateField(c, p.ID, "firstname", u.FirstName)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty FirstName or same FirstName")
		}
		//update user LastName
		if len(u.LastName) != 0 && u.LastName != user.LastName {
			log.Println("Update user LastName! new LastName is:", u.LastName)
			b, err := updateField(c, p.ID, "lastname", u.LastName)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty LastName or same LastName")
		}
		//update user PhoneNumber
		if len(u.PhoneNumber) != 0 && u.PhoneNumber != user.PhoneNumber {
			log.Println("Update user PhoneNumber! new PhoneNumber is:", u.PhoneNumber)
			b, err := updateField(c, p.ID, "phonenumber", u.PhoneNumber)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty PhoneNumber or same PhoneNumber")
		}
		//update user Birthday
		if len(u.Birthday) != 0 && u.Birthday != user.Birthday {
			log.Println("Update user Birthday! new Birthday is:", u.Birthday)
			b, err := updateField(c, p.ID, "birthday", u.Birthday)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty Birthday or same Birthday")
		}
		//update user Gender
		if len(u.Gender) != 0 && u.Gender != user.Gender {
			log.Println("Update user Gender! new Gender is:", u.Gender)
			b, err := updateField(c, p.ID, "gender", u.Gender)
			if !b || err != nil {
				log.Fatalf("UpsertId failed: %s", err)
				responseHandleAccounts(w, r, RspFailed, ReasonOperationFailed, nil)
				return
			}
			responseHandleAccounts(w, r, RspOK, ReasonSuccess, nil)
			return
		} else {
			log.Println("empty Gender or same Gender")
			responseHandleAccounts(w, r, RspOK, ReasonDuplicate, nil)
			return
		}
		//} else {
		//    responseHandleAccounts(w, r, RspFailed, ReasonWrongPw, nil)
		//    return
		//}
		//} else {
		//    responseHandleAccounts(w, r, RspFailed, ReasonMissingParam, nil)
		//    return
		//}
	} else {
		if len(u.Password) == 0 {
			responseHandleAccounts(w, r, http.StatusBadRequest, "Password is empty", nil)
			return
		} else if len(u.AssociatedId) == 0 && len(u.Email) == 0 && len(u.PhoneNumber) == 0 {
			responseHandleAccounts(w, r, http.StatusBadRequest, "AssocitesId/Email/PhoneNumber must be valid at least one", nil)
			return
		}
		u.TimeStamp = time.Now()
		apikey, ok := APIKey(r.Context())
		if ok {
			u.APIKey = apikey
		}
		u.ID = bson.NewObjectId()
		if err := c.Insert(u); err != nil {
			respondErr(w, r, http.StatusInternalServerError, "failed to insert user. err:", err)
			return
		}
		w.Header().Set("Location", "users/"+u.ID.Hex())
		responseHandleAccounts(w, r, http.StatusCreated, ReasonSuccess, nil)
	}
}

func (s *Server) handleAccountsDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
	// Collection User
	c := session.DB("iamhere").C("users")
	u := NewPath(r.URL.Path)
	if !u.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "Cannot delete all users.")
		return
	}
	if err := c.RemoveId(bson.ObjectIdHex(u.ID)); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete user. err:", err)
		return
	}
	respond(w, r, http.StatusOK, nil) // ok
}

func updateField(collection *mgo.Collection, id string, key string, val string) (bool, error) {
	colQuerier := bson.M{"_id": bson.ObjectIdHex(id)}
	change := bson.M{"$set": bson.M{key: val, "timestamp": time.Now()}}
	err := collection.Update(colQuerier, change)
	if err != nil {
		log.Fatalf("UpsertId failed: %s", err)
		return false, err
	}
	return true, nil
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
