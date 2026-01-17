package models

import "time"

type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	AuthorName  string `json:"author_name,omitempty"`
	AuthorEmail string `json:"author_email,omitempty"`
	AuthorImage string `json:"author_image,omitempty"`
}

type CreatePostParams struct {
	UserID  int64
	Title   string
	Content string
}

type UpdatePostParams struct {
	Title   string
	Content string
}
