package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postDetailHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		fmt.Fprintf(w, "Post slug:%s details\n", ps.ByName("slug"))
	}
}
