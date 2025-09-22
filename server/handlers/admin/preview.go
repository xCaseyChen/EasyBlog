package admin

import (
	"easyblog/database"
	"easyblog/handlers/common"
	"easyblog/utils"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func previewHandler(db *gorm.DB) httprouter.Handle {
	const previewTemplateName = "preview"
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		type PostPage struct {
			Title       string
			Category    string
			Tags        []string
			CreatedAt   string
			UpdatedAt   string
			HtmlContent template.HTML
		}

		id := ps.ByName("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			common.RenderInfoPage(w, http.StatusBadRequest, "bad request")
			return
		}
		postBrief, err := gorm.G[database.PostBrief](db).
			Where("id = ?", idInt).
			Take(r.Context())
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
		postDetail, err := gorm.G[database.PostDetail](db).
			Where("id = ?", idInt).
			Take(r.Context())
		if err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Database error: %v", err)
			return
		}

		var htmlContent string
		if htmlContent, err = utils.MarkdownToHTML(postDetail.Content); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to convert markdown to html: %v", err)
			return
		}

		postPage := PostPage{
			Title:       postBrief.Title,
			Category:    postBrief.Category,
			Tags:        postBrief.Tags,
			CreatedAt:   postBrief.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   postBrief.UpdatedAt.UTC().Format(time.RFC3339),
			HtmlContent: template.HTML(htmlContent),
		}
		var htmlString strings.Builder
		if err = tmpl[previewTemplateName].Execute(&htmlString, postPage); err != nil {
			common.RenderInfoPage(w, http.StatusInternalServerError, "internal server error")
			log.Printf("Failed to execute template %s: %v", previewTemplateName, err)
			return
		}

		fmt.Fprint(w, htmlString.String())
	}
}
