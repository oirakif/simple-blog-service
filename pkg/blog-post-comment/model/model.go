package model

import (
	"time"
)

type CreateBlogPostCommentHTTPPayload struct {
	Content string `json:"content" binding:"required"`
}

type BlogPostComment struct {
	ID         *int       `json:"id,omitempty"`
	PostID     *int       `json:"post_id,omitempty"`
	AuthorName *string    `json:"author_name,omitempty"`
	Content    *string    `json:"content,omitempty"`
	Status     *string    `json:"status,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type BlogPostCommentResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type BlogPostCommentPathParam struct {
	ID *int `uri:"id" binding:"required,gt=0"`
}
type GetBlogPostCommentsQueryParams struct {
	ID         *int    `form:"id" binding:"omitempty,gt=0"`
	PostID     *int    `form:"post_id" binding:"omitempty,gt=0"`
	AuthorName *string `form:"author_name" binding:"omitempty"`
	Page       *int    `form:"page" binding:"omitempty,gt=0"`
	PerPage    *int    `form:"per_page" binding:"omitempty,gt=0"`
	Status     *string `form:"status" binding:"omitempty"`
	SortBy     *string `form:"sort_by" binding:"omitempty"`
	SortOrder  *string `form:"sort_order" binding:"omitempty"`
}

type BlogPostCommentFilterQuery struct {
	ID         *int
	PostID     *int
	AuthorName *string
	Status     *string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	SortBy     *string
	SortOrder  *string
	Limit      int
	Offset     int
}
