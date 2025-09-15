package common

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

func setupHandler(db *gorm.DB, jwtSecret string) httprouter.Handle {
	setupTemplateName := "setup"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := gorm.G[database.LocalUser](db).Where("username = ?", "admin").First(r.Context()); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				var htmlString strings.Builder
				err = tmpl[setupTemplateName].Execute(&htmlString, nil)
				if err != nil {
					RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
					log.Printf("Failed to execute template %s: %v", setupTemplateName, err)
					return
				}
				fmt.Fprint(w, htmlString.String())
				return
			} else {
				RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
				log.Printf("Database error: %v", err)
				return
			}
		}
		http.Redirect(w, r, "/home", http.StatusFound)
	}
}
