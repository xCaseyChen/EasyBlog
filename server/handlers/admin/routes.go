package admin

import (
	"easyblog/handlers/common"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

var adminGetHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/admin/dashboard":         dashboardHandler,  // admin dashboard page
	"/admin/posts/:id/preview": previewHandler,    // admin preview post page
	"/api/admin/ping":          pingHandler,       // admin ping api
	"/api/admin/posts/:id":     postsInfoHandler,  // admin post brief and detail api
	"/api/admin/posts":         postsQueryHandler, // admin post list query api
}

var adminPostHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/api/admin/logout": logoutHandler,    // admin logout api
	"/api/admin/posts":  postsPostHandler, // admin posts POST (new post with title)
}

var adminDeleteHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/api/admin/posts/:id": postsDeleteHandler, // admin posts DELETE (delete posts)
}

var adminPatchHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/api/admin/posts/:id": postsPatchHandler, // admin posts PATCH (modify status)
}

var adminPutHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/api/admin/posts/:id": postsPutHandler, // admin posts PUT (update posts)
}

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
	for path, handler := range adminDeleteHandlers {
		r.DELETE(path, withAuth(handler(db), jwtSecret))
	}
	for path, handler := range adminPatchHandlers {
		r.PATCH(path, withAuth(handler(db), jwtSecret))
	}
	for path, handler := range adminPutHandlers {
		r.PUT(path, withAuth(handler(db), jwtSecret))
	}
}
