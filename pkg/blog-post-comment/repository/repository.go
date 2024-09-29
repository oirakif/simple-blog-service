package repository

import (
	"database/sql"
	"log"
	"oirakif/simple-blog-service/pkg/blog-post-comment/model"

	"github.com/Masterminds/squirrel"
)

type BlogPostCommentRepository struct {
	db *sql.DB
}

func NewBlogPostCommentRepository(db *sql.DB) *BlogPostCommentRepository {
	return &BlogPostCommentRepository{
		db: db,
	}
}

func (r *BlogPostCommentRepository) InsertBlogPostComments(newBlogPost model.BlogPostComment) (newID int, err error) {
	builder := squirrel.
		Insert("comments").
		Columns(
			"post_id",
			"author_name",
			"content",
			"status",
			"created_at",
			"updated_at",
		).
		Values(
			*newBlogPost.PostID,
			*newBlogPost.AuthorName,
			*newBlogPost.Content,
			*newBlogPost.Status,
			*newBlogPost.CreatedAt,
			*newBlogPost.UpdatedAt,
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	result, err := r.db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	newIDInt64, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return int(newIDInt64), nil
}
