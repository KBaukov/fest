package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

var NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	log.WithFields(log.Fields{
		"method": r.Method,
		"uri":    r.RequestURI, // r.URL.Path, // r.RequestURI,
		"status": http.StatusNotFound,
	}).Warn("404: Not found")

	Respond(w, Message(false, nil, "Resources was not found"))
})

var NotAllowedMethod = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	log.WithFields(log.Fields{
		"method": r.Method,
		"uri":    r.RequestURI, // r.URL.Path, // r.RequestURI,
		"status": http.StatusMethodNotAllowed,
	}).Warn("405: Not allowed")

	Respond(w, Message(false, nil, "Method not allowed"))
})
