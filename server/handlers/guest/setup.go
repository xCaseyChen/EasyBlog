package guest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"easyblog/database"
)

func setupHandler(db *gorm.DB) httprouter.Handle {
	setupTemplateName := "setup"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").First(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var htmlString strings.Builder
				err = Template.ExecuteTemplate(&htmlString, setupTemplateName, nil)
				if err != nil {
					http.Error(w, "internal server error page\n", http.StatusInternalServerError)
					log.Printf("Failed to execute template %s: %v", setupTemplateName, err)
					return
				}
				fmt.Fprint(w, htmlString.String())
				return
			} else {
				http.Error(w, "Internal server error page\n", http.StatusInternalServerError)
				log.Printf("Database error: %v", err)
				return
			}
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}
