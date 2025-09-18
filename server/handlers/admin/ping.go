package admin

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func pingHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Ping admin")
	}
}
