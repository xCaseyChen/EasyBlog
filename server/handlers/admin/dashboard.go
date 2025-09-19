package admin

import (
	"easyblog/handlers/common"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func dashboardHandler(db *gorm.DB) httprouter.Handle {
	const dashboardTemplateName = "dashboard"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var htmlString strings.Builder
		if err := tmpl[dashboardTemplateName].Execute(&htmlString, nil); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to execute template %s: %v", dashboardTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
