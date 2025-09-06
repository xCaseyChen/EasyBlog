package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func categoriesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Categories page\n")
}
