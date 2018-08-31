package main

import (
	"context"
	"flag"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	lelog "github.com/labstack/gommon/log"
)

var echoInstance = echo.New()

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

	// Echo instance
	//e := echo.New()

	// Middleware
	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Logger.SetLevel(lelog.DEBUG)

	// Routes
	g := echoInstance.Group("/messages")
	g.GET("/", s.handleMessagesGetEcho)
	//g.Use(s.handleMessagesGet)

	// Start server
	echoInstance.Logger.Debug("run echo on 1323 err:")
	echoInstance.Start(":8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/messages/", withCORS(withAPIKey(s.handleMessages)))
	mux.HandleFunc("/areas/", withCORS(withAPIKey(s.handleAreas)))
	mux.HandleFunc("/accounts/", withCORS(withAPIKey(s.handleAccounts)))
	log.Println("Starting web server on", *addr)
	//go http.ListenAndServeTLS(":8082", "../assets/certs/server.crt", "../assets/certs/server.key", mux)
	http.ListenAndServe(":8080", mux)
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
