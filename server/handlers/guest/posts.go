package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func postsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Posts page\n")
}
