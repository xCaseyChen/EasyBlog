package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

func logoutHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		type JsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,  // disable JS read
			Secure:   false, // ! Local test with HTTP -> HTTPS
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Unix(0, 0),
		})
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "logout success",
		})
	}
}
