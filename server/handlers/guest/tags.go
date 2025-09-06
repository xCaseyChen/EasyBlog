package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func tagsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Tags page\n")
}
