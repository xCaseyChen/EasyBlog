package database

import (
	"time"

	"github.com/lib/pq"
)

type PostBrief struct {
	ID        uint           `json:"id"`
	Title     string         `json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Category  string         `json:"category"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags"`
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
