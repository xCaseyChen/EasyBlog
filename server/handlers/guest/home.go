package guest

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"

	"easyblog/database"
	"easyblog/handlers/common"
	"easyblog/utils"
)

func homeHandler(db *gorm.DB) httprouter.Handle {
	homeTemplateName := "home"
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		type HomePage struct {
			HtmlContent template.HTML
		}
		homeBrief, err := gorm.G[database.PostBrief](db).Where("slug = ?", "home").First(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}
		homeDetail, err := gorm.G[database.PostDetail](db).Where("id = ?", homeBrief.ID).First(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}
		var htmlContent string
		if htmlContent, err = utils.MarkdownToHTML(homeDetail.Content); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to convert markdown to html: %v", err)
			return
		}
		homePage := HomePage{
			HtmlContent: template.HTML(htmlContent),
		}
		var htmlString strings.Builder
		if err = tmpl[homeTemplateName].Execute(&htmlString, homePage); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to execute template %s: %v", homeTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
