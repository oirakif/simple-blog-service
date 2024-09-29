package repository

import (
	"database/sql"
	"fmt"
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

func (r *BlogPostCommentRepository) GetBlogPostComments(filterQuery model.BlogPostCommentFilterQuery) (data []model.BlogPostComment, err error) {
	builder := squirrel.
		Select(
			"id",
			"post_id",
			"author_name",
			"content",
			"status",
			"created_at",
			"updated_at").
		From("comments")

	if filterQuery.ID != nil {
		builder = builder.Where("id=?", *filterQuery.ID)
	}

	if filterQuery.PostID != nil {
		builder = builder.Where("post_id=?", *filterQuery.PostID)
	}

	if filterQuery.AuthorName != nil {
		builder = builder.Where("author_name like '%?%'", *filterQuery.AuthorName)
	}

	if filterQuery.Status != nil {
		builder = builder.Where("status=?", *filterQuery.Status)
	}
	builder = builder.Limit(uint64(filterQuery.Limit))
	builder = builder.Offset(uint64(filterQuery.Offset))

	if filterQuery.SortBy != nil {
		if filterQuery.SortOrder != nil {
			builder = builder.OrderBy(fmt.Sprintf("%s %s", *filterQuery.SortBy, *filterQuery.SortOrder))
		} else {
			builder = builder.OrderBy(*filterQuery.SortBy)
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data = []model.BlogPostComment{}
	for rows.Next() {
		var bpc model.BlogPostComment
		err = rows.Scan(
			&bpc.ID,
			&bpc.PostID,
			&bpc.AuthorName,
			&bpc.Content,
			&bpc.Status,
			&bpc.CreatedAt,
			&bpc.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		data = append(data, bpc)
	}

	return data, nil
}
