package common

import (
	"net/http"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	RenderInfoPage(w, http.StatusNotFound, "page not found")
}
