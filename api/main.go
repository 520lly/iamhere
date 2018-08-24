package main

import (
	"context"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", "localhost", "mongodb address")
	)
	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	defer db.Close()
	s := &Server{
		db: db,
	}
	db.SetMode(mgo.Monotonic, true)
	//mux := http.NewServeMux()
	router := mux.NewRouter()
	router.HandleFunc("/authenticate", withCORS(CreateTokenEndpoint)).Methods("POST")
	router.HandleFunc("/messages/", withCORS(ValidateMiddleware(s.handleMessages)))
	router.HandleFunc("/areas/", withCORS(ValidateMiddleware(s.handleAreas)))
	router.HandleFunc("/accounts/", withCORS(ValidateMiddleware(s.handleAccounts)))
	router.HandleFunc("/auth/", withCORS(ValidateMiddleware(s.handleAccounts)))
	log.Println("Starting web server on", *addr)
	//go http.ListenAndServeTLS(":8082", "../assets/certs/server.crt", "../assets/certs/server.key", mux)
	http.ListenAndServe(":8080", router)
	log.Println("Stopping...")
}

// Server is the API server.
type Server struct {
	db *mgo.Session
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

func APIKey(ctx context.Context) (string, bool) {
	key := ctx.Value(contextKeyAPIKey)
	if key == nil {
		return "", false
	}
	keystr, ok := key.(string)
	return keystr, ok
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}
