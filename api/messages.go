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
	ReasonSuccess       string = "Success"
	ReasonFailureParam  string = "Wrong parameter"
	ReasonFailureAPIKey string = "Wrong APIKey"
	ReasonFailueGeneral string = "Failure in general"
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

func (s *Server) handlemessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handlemessagesGet(w, r)
		return
	case "POST":
		log.Println("POST")
		s.handlemessagesPost(w, r)
		return
	case "DELETE":
		s.handlemessagesDelete(w, r)
		return
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
		return
	}
	// not found
	respondHTTPErr(w, r, http.StatusNotFound)
}

func (s *Server) handlemessagesGet(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
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

func (s *Server) handlemessagesPost(w http.ResponseWriter, r *http.Request) {
	log.Println("handlemessagesPost")
	session := s.db.Copy()
	defer session.Close()
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

func (s *Server) handlemessagesDelete(w http.ResponseWriter, r *http.Request) {
	session := s.db.Copy()
	defer session.Close()
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
	}
	result := &response{
		Code:   code,
		Reason: reason,
		Data:   msgs}
	//res, err := json.Marshal(result)
	//if err != nil {
	//    log.Fatalf("JSON marshaling failed: %s", err)
	//}
	//log.Println(string(res))
	respond(w, r, http.StatusOK, &result)
}
