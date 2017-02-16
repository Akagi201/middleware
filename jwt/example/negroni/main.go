package main

import (
	"encoding/json"
	"net/http"

	jwtmiddleware "github.com/Akagi201/middleware/jwt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	startServer()
}

func startServer() {
	r := mux.NewRouter()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("supersecret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	r.HandleFunc("/ping", pingHandler)
	r.Handle("/secured/ping", negroni.New(
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(securedPingHandler)),
	))
	http.Handle("/", r)
	http.ListenAndServe(":3001", nil)
}

type responseText struct {
	Text string `json:"text"`
}

func respondJSON(text string, w http.ResponseWriter) {
	response := responseText{text}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You don't need to be authenticated to call this", w)
}

func securedPingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON("All good. You only get this message if you're authenticated", w)
}
