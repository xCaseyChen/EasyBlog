package guest

import (
	"bytes"
	"easyblog/database"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/julienschmidt/httprouter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"gorm.io/gorm"
)

func postDetailHandler(db *gorm.DB) httprouter.Handle {
	postDetailTemplateName := "post"
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		type Post struct {
			ID          uint
			Title       string
			CreatedAt   time.Time
			UpdatedAt   time.Time
			Category    string
			Tags        []string
			HtmlContent template.HTML
		}
		// get post_brief by slug
		slug := ps.ByName("slug")
		postBrief, err := gorm.G[database.PostBrief](db).Where("slug = ?", slug).First(r.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "page not found\n", http.StatusNotFound)
				return
			} else {
				http.Error(w, "internal server error page\n", http.StatusInternalServerError)
				log.Printf("Database error: %v", err)
				return
			}
		}
		// get post_detail by id
		postDetail, err := gorm.G[database.PostDetail](db).Where("id = ?", postBrief.ID).First(r.Context())
		if err != nil {
			http.Error(w, "internal server error page\n", http.StatusInternalServerError)
			log.Printf("Database error: %v", err)
			return
		}
		// markdown to html
		var htmlContent bytes.Buffer
		highlightGoldmark := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("monokai"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
		)
		if err := highlightGoldmark.Convert([]byte(postDetail.Content), &htmlContent); err != nil {
			http.Error(w, "internal server error page\n", http.StatusInternalServerError)
			log.Printf("Failed to convert markdown to html: %v", err)
			return
		}
		// replace template
		post := Post{
			ID:          postBrief.ID,
			Title:       postBrief.Title,
			CreatedAt:   postBrief.CreatedAt,
			UpdatedAt:   postBrief.UpdatedAt,
			Category:    postBrief.Category,
			Tags:        postBrief.Tags,
			HtmlContent: template.HTML(htmlContent.String()),
		}
		var htmlString strings.Builder
		err = Template.ExecuteTemplate(&htmlString, postDetailTemplateName, post)
		if err != nil {
			http.Error(w, "internal server error page\n", http.StatusInternalServerError)
			log.Printf("Failed to execute template %s: %v", postDetailTemplateName, err)
			return
		}
		fmt.Fprint(w, htmlString.String())
	}
}
