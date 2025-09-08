package common

import (
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page\n", http.StatusNotFound)
}
