package guest

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tags := r.URL.Query()["tags"]
		fmt.Fprintf(w, "Posts query: tags:%v\n", tags)
	}
}
