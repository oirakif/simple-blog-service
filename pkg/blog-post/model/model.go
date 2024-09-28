package model

import (
	"time"
)

type CreateBlogPostHTTPPayload struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type GetBlogPostsQueryParams struct {
	ID        *int    `form:"id" binding:"omitempty,gt=0"`
	AuthorID  *int    `form:"author_id" binding:"omitempty,gt=0"`
	Page      *int    `form:"page" binding:"omitempty,gt=0"`
	PerPage   *int    `form:"per_page" binding:"omitempty,gt=0"`
	Title     *string `form:"title" binding:"omitempty"`
	Status    *string `form:"status" binding:"omitempty"`
	SortBy    *string `form:"sort_by" binding:"omitempty"`
	SortOrder *string `form:"sort_order" binding:"omitempty"`
}

type BlogPost struct {
	ID        int       `json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	AuthorID  int       `json:"author_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type BlogPostResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type BlogPostFilterQuery struct {
	ID        *int
	Title     *string
	Content   *string
	AuthorID  *int
	Status    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	SortBy    *string
	SortOrder *string
	Limit     int
	Offset    int
}
