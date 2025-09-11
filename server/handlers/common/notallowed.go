package common

import (
	"net/http"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	RenderInfoPage(w, http.StatusMethodNotAllowed, "method not allowed")
}
