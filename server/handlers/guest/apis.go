package guest

import (
	"easyblog/database"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
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
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success:    false,
				Message:    "too many categories",
				PostBriefs: nil,
			})
			return
		}
		// use tags and category search in database
		var post_briefs []database.PostBrief
		var query gorm.ChainInterface[database.PostBrief] = gorm.G[database.PostBrief](db)
		if len(tags) > 0 {
			query = query.Where("tags @> ?", pq.Array(tags))
		}
		if len(category) == 1 {
			query = query.Where("category = ?", category[0])
		}
		post_briefs, err := query.Find(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(jsonResponse{
				Success:    false,
				Message:    "internal server error",
				PostBriefs: nil,
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(jsonResponse{
			Success:    true,
			Message:    "find post briefs",
			PostBriefs: post_briefs,
		})
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
			_ = json.NewEncoder(w).Encode(jsonResponse{
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
