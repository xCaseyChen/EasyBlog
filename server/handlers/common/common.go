package common

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes(r *httprouter.Router) {
	r.NotFound = http.HandlerFunc(notFoundHandler)
	r.MethodNotAllowed = http.HandlerFunc(methodNotAllowedHandler)
}
