package common

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

var commonGetHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/setup": setupHandler, // setup page
}

var commonPostHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/api/setup/password": setupPasswordHandler, // setup admin password api
}

func Routes(r *httprouter.Router, db *gorm.DB) {
	r.NotFound = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowed = http.HandlerFunc(methodNotAllowedHandler)
	for path, handler := range commonGetHandlers {
		r.GET(path, handler(db))
	}
	for path, handler := range commonPostHandlers {
		r.POST(path, handler(db))
	}
}
