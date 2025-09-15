package common

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

var commonGetHandlers = map[string]func(*gorm.DB, string) httprouter.Handle{
	"/setup": setupHandler, // setup page
}

var commonPostHandlers = map[string]func(*gorm.DB, string) httprouter.Handle{
	"/api/setup/password": setupPasswordHandler, // setup admin password api
	"/api/login":          adminLoginHandler,    // login as admin api
}

func Routes(r *httprouter.Router, db *gorm.DB, jwtSecret string) {
	r.NotFound = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowed = http.HandlerFunc(methodNotAllowedHandler)
	for path, handler := range commonGetHandlers {
		r.GET(path, handler(db, jwtSecret))
	}
	for path, handler := range commonPostHandlers {
		r.POST(path, handler(db, jwtSecret))
	}
}
