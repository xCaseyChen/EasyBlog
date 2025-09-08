package common

import (
	"net/http"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "405 page\n", http.StatusMethodNotAllowed)
}
