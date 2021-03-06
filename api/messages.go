package main

import (
	"log"
	"net/http"

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
	TimeStamp   int           `json:"timestamp"`
	UserDefAddr string        `json:"userdefaddr" bson:"userdefaddr"`
	ExpiryTime  int           `json:"expirytime"`
	Latitude    int           `json:"latitude"`
	Longitude   int           `json:"longitude"`
	Altitude    int           `json:"altitude"`
	APIKey      string        `json:"apikey"`
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
	// Collection Message
	c := session.DB("iamhere").C("messages")
	var q *mgo.Query
	p := NewPath(r.URL.Path)
	if p.HasID() {
		// get specific message
		q = c.FindId(bson.ObjectIdHex(p.ID))
		//[TODO] more filter for specific messages
	} else {
		// get all messages
		q = c.Find(nil)
	}
	var result []*Message
	if err := q.All(&result); err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	responseHandleMessage(w, r, RspOK, ReasonSuccess, &result)
}

func (s *Server) handleMessagesPost(w http.ResponseWriter, r *http.Request) {
	log.Println("handleMessagesPost")
	session := s.db.Copy()
	defer session.Close()
	// Collection Message
	c := session.DB("iamhere").C("messages")
	var p Message
	if err := decodeBody(r, &p); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read message from request", err)
		return
	}
	apikey, ok := APIKey(r.Context())
	if !ok {
		responseHandleMessage(w, r, RspFailed, ReasonFailureAPIKey, nil)
	}
	p.APIKey = apikey
	p.ID = bson.NewObjectId()
	if err := c.Insert(p); err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to insert message", err)
		return
	}
	w.Header().Set("Location", "messages/"+p.ID.Hex())
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
		Reason string      `json:"reasone"`
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
	//res, err := json.Marshal(result)
	//if err != nil {
	//    log.Fatalf("JSON marshaling failed: %s", err)
	//}
	//log.Println(string(res))
	respond(w, r, http.StatusOK, &result)
}
