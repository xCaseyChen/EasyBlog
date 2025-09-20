package guest

import (
	"easyblog/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"easyblog/utils"
)

func postsQueryHandler(db *gorm.DB) httprouter.Handle {
	const limitMax = 20
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json
		// json response
		type PostQueryData struct {
			PostBriefs   []database.PostBrief `json:"post_briefs"`
			NextBeforeID *uint                `json:"next_before_id"`
		}
		type JsonResponse struct {
			Success bool          `json:"success"`
			Message string        `json:"message"`
			Data    PostQueryData `json:"data"`
		}
		// parse params list
		tags := utils.ParseQueryList(r, "tags", ",")
		category := utils.ParseQueryList(r, "category", ",")
		beforeId := utils.ParseQueryList(r, "before_id", ",")
		limit := utils.ParseQueryList(r, "limit", ",")
		// check category, limit, before_id len
		if len(category) > 1 || len(limit) > 1 || len(beforeId) > 1 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "too many category/limit/before_id",
				Data: PostQueryData{
					PostBriefs:   nil,
					NextBeforeID: nil,
				},
			})
			return
		}
		// use tags, category, before_id, limit search in database
		query := gorm.G[database.PostBrief](db).Where("status = ? ", "published")
		if len(beforeId) == 1 {
			beforeIdInt, err := strconv.Atoi(beforeId[0])
			if err != nil || beforeIdInt <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "invalid before_id: " + beforeId[0],
					Data: PostQueryData{
						PostBriefs:   nil,
						NextBeforeID: nil,
					},
				})
				return
			}
			query = query.Where("id < ?", beforeIdInt)
		}
		if len(tags) > 0 {
			query = query.Where("tags @> ?", pq.Array(tags))
		}
		if len(category) == 1 {
			query = query.Where("category = ?", category[0])
		}
		query = query.Order("pinned DESC, id DESC")
		if len(limit) == 1 {
			limitInt, err := strconv.Atoi(limit[0])
			if err != nil || limitInt <= 0 || limitInt > limitMax {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "invalid limit: " + limit[0],
					Data: PostQueryData{
						PostBriefs:   nil,
						NextBeforeID: nil,
					},
				})
				return
			}
			query = query.Limit(limitInt)
		} else {
			query = query.Limit(limitMax)
		}
		postBriefs, err := query.Find(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
				Data: PostQueryData{
					PostBriefs:   nil,
					NextBeforeID: nil,
				},
			})
			log.Printf("Database error: %v", err)
			return
		}
		var nextBeforeID *uint = nil
		if len(postBriefs) != 0 {
			nextBeforeID = &postBriefs[len(postBriefs)-1].ID
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "find post briefs",
			Data: PostQueryData{
				PostBriefs:   postBriefs,
				NextBeforeID: nextBeforeID,
			},
		})
	}
}

func commentsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tags := r.URL.Query()["tags"]
		fmt.Fprintf(w, "Comments query: tags:%v\n", tags)
	}
}

func allTagsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		type AllTagsData struct {
			Tags pq.StringArray `gorm:"column:all_tags;type:text[]" json:"tags"`
		}
		type JsonResponse struct {
			Success bool         `json:"success"`
			Message string       `json:"message"`
			Data    *AllTagsData `json:"data"`
		}

		var allTagsData AllTagsData
		subQuery := db.Model(&database.PostBrief{}).Select("unnest(tags) AS tag")
		if err := db.Table("(?) as t", subQuery).Select("array_agg(DISTINCT tag ORDER BY tag) AS all_tags").Scan(&allTagsData).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			Message: "find all tags",
			Data:    &allTagsData,
		})
	}
}

func allCategoriesQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")

		type AllCategoriesData struct {
			Categories pq.StringArray `gorm:"column:all_categories;type:text[]" json:"categories"`
		}
		type JsonResponse struct {
			Success bool               `json:"success"`
			Message string             `json:"message"`
			Data    *AllCategoriesData `json:"data"`
		}

		var allCategoriesData AllCategoriesData
		if err := db.Model(&database.PostBrief{}).Where("category IS NOT NULL").Select("array_agg(DISTINCT category ORDER BY category) AS all_categories").Scan(&allCategoriesData).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
			Message: "find all categories",
			Data:    &allCategoriesData,
		})
	}
}
