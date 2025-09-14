package guest

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

// <URL, Handle> maps
var guestGetHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/":                   rootHandler,          // root page (redirect)
	"/home":               homeHandler,          // home page
	"/list":               listHandler,          // posts page
	"/categories":         categoriesHandler,    // categories page
	"/tags":               tagsHandler,          // tags page
	"/about":              aboutHandler,         // about page
	"/post/:slug":         postDetailHandler,    // post detail page by slug
	"/api/guest/posts":    postsQueryHandler,    // query posts api
	"/api/guest/comments": commentsQueryHandler, // query comments api
}

var guestPostHandlers = map[string]func(*gorm.DB) httprouter.Handle{}

// Add guest handlers to router
func Routes(r *httprouter.Router, db *gorm.DB) {
	for path, handler := range guestGetHandlers {
		r.GET(path, handler(db))
	}
	for path, handler := range guestPostHandlers {
		r.POST(path, handler(db))
	}
}
