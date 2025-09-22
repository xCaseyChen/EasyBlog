package admin

import (
	"easyblog/database"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsInfoHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json

		type InfoData struct {
			Title    string   `json:"title"`
			Category string   `json:"category"`
			Tags     []string `json:"tags"`
			Status   string   `json:"status"`
			Pinned   bool     `json:"pinned"`
			Content  string   `json:"content"`
		}

		type JsonResponse struct {
			Success bool      `json:"success"`
			Message string    `json:"message"`
			Data    *InfoData `json:"data"`
		}

		id := ps.ByName("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request: id",
				Data:    nil,
			})
			return
		}

		postBrief, err := gorm.G[database.PostBrief](db).
			Where("id = ?", idInt).
			Take(r.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "post not found",
					Data:    nil,
				})
				return
			} else {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "internal server error",
					Data:    nil,
				})
				log.Printf("Database error: %v", err)
				return
			}
		}

		postDetail, err := gorm.G[database.PostDetail](db).
			Where("id = ?", idInt).
			Take(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
				Data:    nil,
			})
			log.Printf("Database error: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "fetch post info",
			Data: &InfoData{
				Title:    postBrief.Title,
				Category: postBrief.Category,
				Tags:     postBrief.Tags,
				Status:   postBrief.Status,
				Pinned:   postBrief.Pinned,
				Content:  postDetail.Content,
			},
		})
	}
}
