package admin

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsPatchHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)     // limit body to 1 MB
		w.Header().Set("Content-Type", "application/json") // response type: json

		type JsonRequest struct {
			Status string `json:"status"`
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
		if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
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

		rowsAffected, err := gorm.G[database.PostBrief](db).
			Where("id = ?", idInt).
			Update(r.Context(), "status", jsonRequest.Status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Database error: %v", err)
			return
		}
		if rowsAffected == 0 {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "post not found",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "status update",
		})
	}
}
