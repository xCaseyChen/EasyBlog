package guest

import (
	"easyblog/handlers/common"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func listHandler(db *gorm.DB) httprouter.Handle {
	const listTemplateName = "list"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var htmlString strings.Builder
		if err := tmpl[listTemplateName].Execute(&htmlString, nil); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to execute template %s: %v", listTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
