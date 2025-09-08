package database

import (
	"time"
)

type PostBrief struct {
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Category  string
	Tags      []string
}

type PostDetail struct {
	ID      uint
	Content string
}

type Comment struct {
	ID        uint
	PostID    uint
	CreatedAt time.Time
	Content   string
}

type LocalUser struct {
	Username string
	Password string
}
