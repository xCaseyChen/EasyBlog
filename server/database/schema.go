package database

import (
	"time"

	"github.com/lib/pq"
)

type PostBrief struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	Slug      string         `json:"slug"`
	Category  string         `json:"category"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
	Status    string         `json:"status"`
	Pinned    bool           `json:"pinned"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type PostDetail struct {
	ID      uint
	Content string
}

type Comment struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type LocalUser struct {
	Username string
	Password string
}
