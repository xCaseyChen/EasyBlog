package guest

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

// <URL, Handle> maps
var guestGetHandlers = map[string]func(*gorm.DB) httprouter.Handle{
	"/":                     rootHandler,               // root page (redirect)
	"/home":                 homeHandler,               // home page
	"/list":                 listHandler,               // posts page
	"/categories":           categoriesHandler,         // categories page
	"/tags":                 tagsHandler,               // tags page
	"/about":                aboutHandler,              // about page
	"/posts/:slug":          postDetailHandler,         // post detail page by slug
	"/api/guest/posts":      postsQueryHandler,         // query posts api
	"/api/guest/comments":   commentsQueryHandler,      // query comments api
	"/api/guest/tags":       allTagsQueryHandler,       // query all tags
	"/api/guest/categories": allCategoriesQueryHandler, // query all categories
}

// Add guest handlers to router
func Routes(r *httprouter.Router, db *gorm.DB) {
	for path, handler := range guestGetHandlers {
		r.GET(path, handler(db))
	}
}
