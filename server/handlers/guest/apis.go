package guest

import (
	"easyblog/database"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"easyblog/utils"
)

func postsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json
		// json response
		type jsonResponse struct {
			Success    bool                 `json:"success"`
			Message    string               `json:"message"`
			PostBriefs []database.PostBrief `json:"post_briefs"`
		}
		// parse params list
		tags := utils.ParseQueryList(r, "tags", ",")
		category := utils.ParseQueryList(r, "category", ",")
		// check category
		if len(category) > 1 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(jsonResponse{false, "to many categories", nil})
			return
		}
		// use tags and category search in database
		log.Panicln(tags)
	}
}

func commentsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tags := r.URL.Query()["tags"]
		fmt.Fprintf(w, "Comments query: tags:%v\n", tags)
	}
}

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
			json.NewEncoder(w).Encode(jsonResponse{
				Success: false,
				Message: "invalid request",
			})
			log.Printf("Failed to parse http.Request.Body: %v", err)
			return
		}
		// TODO: validate password, uppper case & lower case & number & len > 8
		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(jsonReq.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(jsonResponse{
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
				json.NewEncoder(w).Encode(jsonResponse{
					Success: false,
					Message: "admin already exists",
				})
				log.Printf("Admin already exists: %v", err)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(jsonResponse{
					Success: false,
					Message: "internal server error",
				})
				log.Printf("Database error: %v", err)
				return
			}
		} else {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(jsonResponse{
				Success: true,
				Message: "admin password set up",
			})
			log.Printf("Admin password set up")
			return
		}
	}
}
