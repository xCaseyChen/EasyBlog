package admin

import (
	"easyblog/database"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func postsDeleteHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json

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
				Message: "invalid request",
			})
			return
		}

		transaction := func(tx *gorm.DB) error {
			postBrief, err := gorm.G[database.PostBrief](tx).
				Where("id = ?", idInt).
				Take(r.Context())
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("post not found")
				}
				return err
			}
			if postBrief.Status != "deleted" {
				return fmt.Errorf("post not in deleted status")
			}
			_, err = gorm.G[database.PostBrief](tx).
				Where("id = ?", idInt).
				Delete(r.Context())
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
			log.Printf("Failed to delete post (database transaction): %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "post deleted",
		})
	}
}
