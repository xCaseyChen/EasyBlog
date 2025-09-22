package admin

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsPutHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 5<<20)     // ! limit body to 5 MB
		w.Header().Set("Content-Type", "application/json") // response type: json

		type JsonRequest struct {
			Title    string   `json:"title"`
			Category string   `json:"category"`
			Tags     []string `json:"tags"`
			Status   string   `json:"status"`
			Pinned   bool     `json:"pinned"`
			Content  string   `json:"content"`
		}
		type JsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		id := ps.ByName("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request: id",
			})
			return
		}

		var jsonRequest JsonRequest
		if err = json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request: json",
			})
			return
		}
		if !utils.IsValidStatus(jsonRequest.Status) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request: status",
			})
			return
		}

		transaction := func(tx *gorm.DB) error {
			_, err = gorm.G[database.PostBrief](tx).Where("id = ?", idInt).Updates(r.Context(), database.PostBrief{
				Title:    jsonRequest.Title,
				Slug:     fmt.Sprintf("%s-%d", jsonRequest.Title, idInt),
				Category: jsonRequest.Category,
				Tags:     jsonRequest.Tags,
				Status:   jsonRequest.Status,
				Pinned:   jsonRequest.Pinned,
			})
			if err != nil {
				return err
			}
			_, err = gorm.G[database.PostDetail](tx).Where("id = ?", idInt).Update(r.Context(), "content", jsonRequest.Content)
			if err != nil {
				return err
			}
			return nil
		}

		if err = db.Transaction(transaction); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Failed to update post (database transaction): %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "update post",
		})
	}
}
