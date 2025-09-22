package admin

import (
	"easyblog/database"
	"easyblog/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func postsQueryHandler(db *gorm.DB) httprouter.Handle {
	const limitMax = 10
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json

		type QueryData struct {
			PostBriefs []database.PostBrief `json:"post_briefs"`
			Page       int                  `json:"page"`
			Limit      int                  `json:"limit"`
			TotalPages int                  `json:"total_pages"`
		}
		type JsonResponse struct {
			Success bool       `json:"success"`
			Message string     `json:"message"`
			Data    *QueryData `json:"data"`
		}

		var err error
		page := utils.ParseQueryList(r, "page", ",")
		limit := utils.ParseQueryList(r, "limit", ",")
		category := utils.ParseQueryList(r, "category", ",")
		tags := utils.ParseQueryList(r, "tags", ",")
		status := utils.ParseQueryList(r, "status", ",")
		pinned := utils.ParseQueryList(r, "pinned", ",")
		title_like := utils.ParseQueryList(r, "title_like", ",")

		if len(page) > 1 ||
			len(limit) > 1 ||
			len(category) > 1 ||
			len(status) != 1 ||
			len(pinned) > 1 ||
			len(title_like) > 1 {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid query parameters",
				Data:    nil,
			})
			return
		}

		pageInt := 1
		if len(page) == 1 {
			pageInt, err = strconv.Atoi(page[0])
			if err != nil || pageInt <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "invalid query parameter: page",
					Data:    nil,
				})
				return
			}
		}
		limitInt := limitMax
		if len(limit) == 1 {
			limitInt, err = strconv.Atoi(limit[0])
			if err != nil || limitInt <= 0 || limitInt > limitMax {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "invalid query parameter: limit",
					Data:    nil,
				})
				return
			}
		}
		query := gorm.G[database.PostBrief](db).
			Limit(limitInt).
			Offset((pageInt - 1) * limitInt)
		if len(category) == 1 {
			query = query.Where("category = ?", category[0])
		}
		if len(tags) >= 1 {
			query = query.Where("tags @> ?", pq.Array(tags))
		}
		if !utils.IsValidStatus(status[0]) {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "invalid query parameter: status",
				Data:    nil,
			})
			return
		} else {
			query = query.Where("status = ?", status[0])
		}
		if len(pinned) == 1 {
			pinnedBool, err := strconv.ParseBool(pinned[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success: false,
					Message: "invalid query parameter: pinned",
					Data:    nil,
				})
				return
			}
			query = query.Where("pinned = ?", pinnedBool)
		}
		if len(title_like) == 1 {
			query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", title_like[0]))
		}
		query = query.Order("id DESC")
		postBriefs, err := query.Find(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
				Data:    nil,
			})
			log.Printf("Database error: %v", err)
			return
		}
		var totalPostCount int64
		if err = db.Model(&database.PostBrief{}).Count(&totalPostCount).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success: false,
				Message: "internal server error",
				Data:    nil,
			})
			log.Printf("Database error: %v", err)
			return
		}
		totalPages := int(math.Ceil(float64(totalPostCount) / float64(limitInt)))

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success: true,
			Message: "fetch post list",
			Data: &QueryData{
				PostBriefs: postBriefs,
				Page:       pageInt,
				Limit:      limitInt,
				TotalPages: totalPages,
			},
		})
	}
}
