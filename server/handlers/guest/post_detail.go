package guest

import (
	"easyblog/database"
	"easyblog/handlers/common"
	"easyblog/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"html/template"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postDetailHandler(db *gorm.DB) httprouter.Handle {
	postDetailTemplateName := "post"
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		type PostPage struct {
			Title       string
			CreatedAt   string
			UpdatedAt   string
			Category    string
			Tags        []string
			HtmlContent template.HTML
		}
		// get post_brief by slug
		slug := ps.ByName("slug")
		postBrief, err := gorm.G[database.PostBrief](db).Where("slug = ?", slug).First(r.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				common.RenderInfoPage(w, http.StatusNotFound, "page not found")
				return
			} else {
				common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
				log.Printf("Database error: %v", err)
				return
			}
		}
		// get post_detail by id
		postDetail, err := gorm.G[database.PostDetail](db).Where("id = ?", postBrief.ID).First(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}
		// markdown to html
		var htmlContent string
		if htmlContent, err = utils.MarkdownToHTML(postDetail.Content); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to convert markdown to html: %v", err)
			return
		}
		// replace template
		postPage := PostPage{
			Title:       postBrief.Title,
			CreatedAt:   postBrief.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   postBrief.UpdatedAt.UTC().Format(time.RFC3339),
			Category:    postBrief.Category,
			Tags:        postBrief.Tags,
			HtmlContent: template.HTML(htmlContent),
		}
		var htmlString strings.Builder
		err = tmpl.ExecuteTemplate(&htmlString, postDetailTemplateName, postPage)
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to execute template %s: %v", postDetailTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
