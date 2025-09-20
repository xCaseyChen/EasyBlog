package common

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupPasswordHandler(db *gorm.DB, jwtSecret string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)     // limit body to 1 MB
		w.Header().Set("Content-Type", "application/json") // response type: json
		// json request and response
		type JsonRequest struct {
			Password string `json:"password"`
		}
		type JsonResponse struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		}
		// get password
		var jsonRequest JsonRequest
		if err := json.NewDecoder(r.Body).Decode(&jsonRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid request",
			})
			log.Printf("Failed to parse http.Request.Body: %v", err)
			return
		}
		// validate password compliance
		if !utils.ValidatePasswordCompliance(jsonRequest.Password) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "incompliance password",
			})
			return
		}
		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(jsonRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
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
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "admin already exists",
				})
				log.Printf("Admin already exists: %v", err)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "internal server error",
				})
				log.Printf("Database error: %v", err)
				return
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: true,
				Message: "admin password set up",
			})
			log.Printf("Admin password set up")
			return
		}
	}
}

func adminLoginHandler(db *gorm.DB, jwtSecret string) httprouter.Handle {
	const expire = 24 * time.Hour
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)     // limit body to 1 MB
		w.Header().Set("Content-Type", "application/json") // response type: json
		type JsonRequest struct {
			Password string `json:"password"`
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

		localUser, err := gorm.G[database.LocalUser](db).
			Where("username = ?", "admin").
			Take(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Database error: %v", err)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(localUser.Password), []byte(jsonRequest.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid password",
			})
			return
		}

		tokenString, err := utils.SignJWT([]byte(jwtSecret), jwt.RegisteredClaims{
			Subject:   "admin",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
			})
			log.Printf("Failed to sign JWT: %v", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,  // disable JS read
			Secure:   false, // ! Local test with HTTP -> HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   int(expire.Seconds()),
		})
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "login success",
		})
	}
}
