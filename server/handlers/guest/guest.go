package guest

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

// <URL, Handle> map, method is fixed to "GET"
var guestHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/":             homeHandler,          // home page
	"/home":         homeHandler,          // home page
	"/posts":        postsHandler,         // posts page
	"/categories":   categoriesHandler,    // categories page
	"/tags":         tagsHandler,          // tags page
	"/about":        aboutHandler,         // about page
	"/post/:id":     postDetailHandler,    // post detail page
	"/api/posts":    postsQueryHandler,    // query posts api
	"/api/comments": commentsQueryHandler, // query comments api
}

// Add guest handlers to router
func Routes(r *httprouter.Router, db *gorm.DB) {
	for path, handler := range guestHandlers {
		r.GET(path, handler(db))
	}
}
