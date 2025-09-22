package admin

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsPostHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)     // limit body to 1 MB
		w.Header().Set("Content-Type", "application/json") // response type: json

		type JsonRequest struct {
			Title string `json:"title"`
		}
		type JsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		var jsonRequest JsonRequest
		if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request",
			})
			return
		}

		transaction := func(tx *gorm.DB) error {
			// create post brief with title
			postBrief := database.PostBrief{
				Title:  jsonRequest.Title,
				Status: "draft",
			}
			if err := gorm.G[database.PostBrief](tx).Create(r.Context(), &postBrief); err != nil {
				return err
			}
			// update post slug
			slug := fmt.Sprintf("%s-%d", utils.Slugify(postBrief.Title), postBrief.ID)
			_, err := gorm.G[database.PostBrief](tx).
				Where("id = ?", postBrief.ID).
				Update(r.Context(), "slug", slug)
			if err != nil {
				return err
			}
			// create post detail
			postDetail := database.PostDetail{
				ID:      postBrief.ID,
				Content: "",
			}
			if err := gorm.G[database.PostDetail](tx).Create(r.Context(), &postDetail); err != nil {
				return err
			}
			return nil
		}
		if err := db.Transaction(transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Failed to create post (database transaction): %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "create post",
		})
	}
}
