package guest

import (
	"easyblog/database"
	"easyblog/handlers/common"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func rootHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").Take(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Redirect(w, r, "/setup", http.StatusFound)
			} else {
				common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
				log.Printf("Database error: %v", err)
			}
			return
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}
