package guest

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"easyblog/database"
)

func setupHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").First(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Fprint(w, "Setup page\n")
			} else {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error page\n", http.StatusInternalServerError)
			}
		} else {
			http.Redirect(w, r, "/home", http.StatusFound)
		}
	}
}
