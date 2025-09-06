package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func commentsQueryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tags := r.URL.Query()["tags"]
	fmt.Fprintf(w, "Comments query: tags:%v\n", tags)
}
