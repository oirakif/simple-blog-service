package repository

import (
	"database/sql"
	"fmt"
	"log"
	"oirakif/simple-blog-service/pkg/blog-post/model"

	"github.com/Masterminds/squirrel"
)

type BlogPostRepository struct {
	db *sql.DB
}

func NewBlogPostRepository(db *sql.DB) *BlogPostRepository {
	return &BlogPostRepository{
		db: db,
	}
}

func (r *BlogPostRepository) InsertBlogPost(newBlogPost model.BlogPost) (newID int, err error) {
	builder := squirrel.
		Insert("posts").
		Columns(
			"title",
			"content",
			"author_id",
			"status",
			"created_at",
			"updated_at",
		).
		Values(
			*newBlogPost.Title,
			*newBlogPost.Content,
			*newBlogPost.AuthorID,
			*newBlogPost.Status,
			*newBlogPost.CreatedAt,
			*newBlogPost.UpdatedAt,
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return 0, err
	}

	result, err := r.db.Exec(sql, args...) // ? = placeholder
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

func (r *BlogPostRepository) GetBlogPosts(filterQuery model.BlogPostFilterQuery) (data []model.BlogPost, err error) {
	builder := squirrel.
		Select(
			"id",
			"title",
			"content",
			"author_id",
			"status",
			"created_at",
			"updated_at").
		From("posts")

	if filterQuery.ID != nil {
		builder = builder.Where("id=?", *filterQuery.ID)
	}

	if filterQuery.Title != nil {
		builder = builder.Where("title like '%?%'", *filterQuery.Title)
	}

	if filterQuery.AuthorID != nil {
		builder = builder.Where("author_id=?", *filterQuery.AuthorID)
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

	data = []model.BlogPost{}
	for rows.Next() {
		var bp model.BlogPost
		err = rows.Scan(
			&bp.ID,
			&bp.Title,
			&bp.Content,
			&bp.AuthorID,
			&bp.Status,
			&bp.CreatedAt,
			&bp.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		data = append(data, bp)
	}

	return data, nil
}

func (r *BlogPostRepository) UpdateBlogPost(filterQuery model.BlogPostFilterQuery, updatePayload model.BlogPost) (err error) {
	builder := squirrel.
		Update("posts")

	if updatePayload.Title != nil {
		builder = builder.Set("title", *updatePayload.Title)
	}
	if updatePayload.Content != nil {
		builder = builder.Set("content", *updatePayload.Content)
	}
	if updatePayload.Status != nil {
		builder = builder.Set("status", *updatePayload.Status)
	}
	if updatePayload.UpdatedAt != nil {
		builder = builder.Set("updated_at", *updatePayload.UpdatedAt)
	}

	if filterQuery.ID != nil {
		builder = builder.Where("id=?", *filterQuery.ID)
	}
	if filterQuery.AuthorID != nil {
		builder = builder.Where("author_id=?", *filterQuery.AuthorID)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = r.db.Exec(sql, args...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
