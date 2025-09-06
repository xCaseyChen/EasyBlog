package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func postsQueryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tags := r.URL.Query()["tags"]
	fmt.Fprintf(w, "Posts query: tags:%v\n", tags)
}
