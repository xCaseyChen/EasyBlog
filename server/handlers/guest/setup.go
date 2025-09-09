package guest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"easyblog/database"
)

func setupHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").First(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				htmlBytes, err := os.ReadFile("templates/setup/index.html")
				if err != nil {
					log.Printf("Template setup read error: %v", err)
					http.Error(w, "Internal server error page\n", http.StatusInternalServerError)
					return
				} else {
					fmt.Fprint(w, string(htmlBytes))
					return
				}
			} else {
				log.Printf("Database error: %v", err)
				http.Error(w, "Internal server error page\n", http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/home", http.StatusFound)
			return
		}
	}
}
