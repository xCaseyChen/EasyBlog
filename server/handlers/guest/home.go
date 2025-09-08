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

func homeHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").First(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Redirect(w, r, "/setup", http.StatusFound)
			} else {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error page\n", http.StatusInternalServerError)
			}
		} else {
			fmt.Fprint(w, "Home page\n")
		}
	}
}
