package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func postDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "Post id:%s details\n", ps.ByName("id"))
}
