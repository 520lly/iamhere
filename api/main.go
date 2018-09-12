package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/messages/", withCORS(withAPIKey(s.handleMessages)))
	mux.HandleFunc("/areas/", withCORS(withAPIKey(s.handleAreas)))
	mux.HandleFunc("/accounts/", withCORS(withAPIKey(s.handleAccounts)))

	var m *autocert.Manager
	var httpsSrv *http.Server

	hostPolicy := func(ctx context.Context, host string) error {
		// Note: change to your real host
		allowedHost := "www.historystest.com"
		if host == allowedHost {
			return nil
		}
		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	dataDir := "/home/jaycee/var/www/cache"
	m = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(dataDir),
	}

	httpsSrv = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	httpsSrv.Addr = ":443"
	httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	go func() {
		log.Println("Starting HTTPS server on %s\n", httpsSrv.Addr)
		err := httpsSrv.ListenAndServeTLS("", "")
		if err != nil {
			log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
		}
	}()
	//cfg := &tls.Config{
	//    MinVersion:               tls.VersionTLS12,
	//    CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	//    PreferServerCipherSuites: true,
	//    CipherSuites: []uint16{
	//        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	//        tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	//        tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	//        tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	//    },
	//}
	//srv := &http.Server{
	//    Addr:         ":443",
	//    Handler:      mux,
	//    TLSConfig:    cfg,
	//    TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	//}

	//log.Println("Starting https web server on: :443")
	//go srv.ListenAndServeTLS("/etc/ssl/iamhere/server.crt", "/etc/ssl/iamhere/server.key")
	//go srv.ListenAndServeTLS("../assets/214987401110045.pem", "../assets/214987401110045.key")
	log.Println("Starting web server on", *addr)
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
