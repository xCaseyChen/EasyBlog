package guest

import (
	"easyblog/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"easyblog/utils"
)

func postsQueryHandler(db *gorm.DB) httprouter.Handle {
	limitMax := 20
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json") // response type: json
		// json response
		type JsonResponse struct {
			Success      bool                 `json:"success"`
			Message      string               `json:"message"`
			PostBriefs   []database.PostBrief `json:"post_briefs"`
			NextBeforeID *uint                `json:"next_before_id"`
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
				Success:      false,
				Message:      "too many category/limit/before_id",
				PostBriefs:   nil,
				NextBeforeID: nil,
			})
			return
		}
		// use tags, category, before_id, limit search in database
		query := gorm.G[database.PostBrief](db).Order("id desc")
		if len(limit) == 1 {
			limitInt, err := strconv.Atoi(limit[0])
			if err != nil || limitInt <= 0 || limitInt > limitMax {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success:      false,
					Message:      "invalid limit: " + limit[0],
					PostBriefs:   nil,
					NextBeforeID: nil,
				})
				return
			}
			query = query.Limit(limitInt)
		} else {
			query = query.Limit(limitMax)
		}
		if len(beforeId) == 1 {
			beforeIdInt, err := strconv.Atoi(beforeId[0])
			if err != nil || beforeIdInt <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(JsonResponse{
					Success:      false,
					Message:      "invalid before_id: " + beforeId[0],
					PostBriefs:   nil,
					NextBeforeID: nil,
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
		postBriefs, err := query.Find(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JsonResponse{
				Success:      false,
				Message:      "internal server error",
				PostBriefs:   nil,
				NextBeforeID: nil,
			})
			return
		}
		var nextBeforeID *uint = nil
		if len(postBriefs) != 0 {
			nextBeforeID = &postBriefs[len(postBriefs)-1].ID
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(JsonResponse{
			Success:      true,
			Message:      "find post briefs",
			PostBriefs:   postBriefs,
			NextBeforeID: nextBeforeID,
		})
	}
}

func commentsQueryHandler(db *gorm.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		tags := r.URL.Query()["tags"]
		fmt.Fprintf(w, "Comments query: tags:%v\n", tags)
	}
}
