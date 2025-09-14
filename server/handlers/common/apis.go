package common

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupPasswordHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)     // limit body to 1 MB
		w.Header().Set("Content-Type", "application/json") // response type: json
		// json request and response
		type jsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		type jsonRequest struct {
			Password string `json:"password"`
		}
		// get password
		var jsonReq jsonRequest
		if err := json.NewDecoder(r.Body).Decode(&jsonReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success: false,
				Message: "invalid request",
			})
			log.Printf("Failed to parse http.Request.Body: %v", err)
			return
		}
		// validate password compliance
		if !utils.ValidatePasswordCompliance(jsonReq.Password) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success: false,
				Message: "incompliance password",
			})
			return
		}
		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(jsonReq.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Failed to generate hash: %v", err)
			return
		}
		// write to database
		user := database.LocalUser{
			Username: "admin",
			Password: string(hash),
		}
		err = gorm.G[database.LocalUser](db).Create(r.Context(), &user)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(jsonResponse{
					Success: false,
					Message: "admin already exists",
				})
				log.Printf("Admin already exists: %v", err)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(jsonResponse{
					Success: false,
					Message: "internal server error",
				})
				log.Printf("Database error: %v", err)
				return
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success: true,
				Message: "admin password set up",
			})
			log.Printf("Admin password set up")
			return
		}
	}
}
