package admin

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func pingHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		type JsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "ping success",
		})
	}
}
