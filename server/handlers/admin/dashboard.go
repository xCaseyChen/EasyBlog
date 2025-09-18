package admin

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func dashboardHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Admin Dashboard page\n")
	}
}
