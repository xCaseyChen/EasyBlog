package guest

import (
	"easyblog/database"
	"easyblog/handlers/common"
	"easyblog/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func aboutHandler(db *gorm.DB) httprouter.Handle {
	const aboutTemplateName = "about"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		type AboutPage struct {
			HtmlContent template.HTML
		}
		aboutBrief, err := gorm.G[database.PostBrief](db).Where("slug = ?", "about").Take(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}
		aboutDetail, err := gorm.G[database.PostDetail](db).Where("id = ?", aboutBrief.ID).Take(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}
		var htmlContent string
		if htmlContent, err = utils.MarkdownToHTML(aboutDetail.Content); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to convert markdown to html: %v", err)
			return
		}
		aboutPage := AboutPage{
			HtmlContent: template.HTML(htmlContent),
		}
		var htmlString strings.Builder
		if err = tmpl[aboutTemplateName].Execute(&htmlString, aboutPage); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed execute template %s: %v", aboutTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
