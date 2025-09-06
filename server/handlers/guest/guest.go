package guest

import (
	"github.com/julienschmidt/httprouter"
)

// <URL, Handle> map, method is fixed to "GET"
var guestHandlers = map[string]httprouter.Handle{
	"/":           homeHandler,       // home page
	"/home":       homeHandler,       // home page
	"/posts":      postsHandler,      // posts page
	"/categories": categoriesHandler, // categories page
	"/tags":       tagsHandler,       // tags page
	"/about":      aboutHandler,      // about page
	"/post/:id":   postDetailHandler, // post detail page
	"/api/posts":  postQueryHandler,  // query posts api
}

// Add guest handlers to router
func Routes(r *httprouter.Router) {
	for path, handler := range guestHandlers {
		r.GET(path, handler)
	}
}
