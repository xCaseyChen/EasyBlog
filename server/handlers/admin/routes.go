package admin

import (
	"easyblog/handlers/common"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

var adminGetHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/admin/dashboard": dashboardHandler, // admin dashboard page
}

var adminPostHandlers = map[string]func(*gorm.DB) httprouter.Handle{}

func withAuth(h httprouter.Handle, jwtSecret string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if err := auth(r, jwtSecret); err != nil {
			common.RenderInfoPage(w, http.StatusUnauthorized, "access denied")
			log.Printf("Unauthorized access attempt: %v", err)
			return
		}
		h(w, r, ps)
	}
}

func Routes(r *httprouter.Router, db *gorm.DB, jwtSecret string) {
	for path, handler := range adminGetHandlers {
		r.GET(path, withAuth(handler(db), jwtSecret))
	}
	for path, handler := range adminPostHandlers {
		r.POST(path, withAuth(handler(db), jwtSecret))
	}
}
