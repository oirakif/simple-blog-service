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
	PostID *int `uri:"postID" binding:"required,gt=0"`
}
