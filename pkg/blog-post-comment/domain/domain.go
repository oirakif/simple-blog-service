package domain

import (
	"net/http"
	"oirakif/simple-blog-service/pkg/blog-post-comment/model"
	"oirakif/simple-blog-service/pkg/blog-post-comment/repository"
	blogPostModel "oirakif/simple-blog-service/pkg/blog-post/model"
	blogPostRepository "oirakif/simple-blog-service/pkg/blog-post/repository"
	"oirakif/simple-blog-service/pkg/utils"

	"time"
)

type BlogPostCommentDomain struct {
	blogPostCommentRepository *repository.BlogPostCommentRepository
	blogPostRepository        *blogPostRepository.BlogPostRepository
}

func NewBlogPostCommentDomain(
	blogPostCommentRepository *repository.BlogPostCommentRepository,
	blogPostRepository *blogPostRepository.BlogPostRepository,
) *BlogPostCommentDomain {

	return &BlogPostCommentDomain{
		blogPostCommentRepository: blogPostCommentRepository,
		blogPostRepository:        blogPostRepository,
	}
}

func (d *BlogPostCommentDomain) CreateBlogPostComment(authorName string, postID *int, content string) (statusCode int, response model.BlogPostCommentResponse) {
	status := "ACTIVE"
	filterQuery := blogPostModel.BlogPostFilterQuery{
		ID:     postID,
		Status: &status,
		Limit:  1,
		Offset: 0,
	}

	blogPosts, err := d.blogPostRepository.GetBlogPosts(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog posts data"

		return http.StatusInternalServerError, response
	}
	if len(blogPosts) == 0 {
		response.Error = true
		response.Message = "blog post is not found"

		return http.StatusNotFound, response
	}

	currentTimestamp := time.Now()
	newBlogPostComment := model.BlogPostComment{
		PostID:     postID,
		AuthorName: &authorName,
		Content:    &content,
		Status:     &status,
		CreatedAt:  &currentTimestamp,
		UpdatedAt:  &currentTimestamp,
	}

	newBlogPostID, err := d.blogPostCommentRepository.InsertBlogPostComments(newBlogPostComment)
	if err != nil {
		response.Error = true
		response.Message = "error occured while creating new blog post comment"

		return http.StatusInternalServerError, response
	}

	newBlogPostComment.ID = &newBlogPostID

	response.Data = newBlogPostComment
	response.Message = "new blog post comment created"

	return http.StatusCreated, response
}

func (d *BlogPostCommentDomain) GetBlogPostComments(queries *model.GetBlogPostCommentsQueryParams) (statusCode int, response model.BlogPostCommentResponse) {
	page, perPage := 1, 25
	sortBy, sortOrder := "created_at", "desc"
	status := "ACTIVE"
	var filterQuery model.BlogPostCommentFilterQuery
	if queries.Page != nil {
		page = *queries.Page
	}
	if queries.PerPage != nil {
		perPage = *queries.PerPage
	}
	if queries.SortBy != nil {
		sortBy = *queries.SortBy
	}
	if queries.SortOrder != nil {
		sortOrder = *queries.SortOrder
	}
	if queries.Status != nil {
		status = *queries.Status
	}
	offset := utils.CalculateOffset(page, perPage)

	filterQuery.ID = queries.ID
	filterQuery.PostID = queries.PostID
	filterQuery.AuthorName = queries.AuthorName
	filterQuery.Status = &status
	filterQuery.Limit = perPage
	filterQuery.Offset = offset
	filterQuery.SortBy = &sortBy
	filterQuery.SortOrder = &sortOrder

	blogPosts, err := d.blogPostCommentRepository.GetBlogPostComments(filterQuery)
	if err != nil {
		response.Error = true
		response.Message = "error occured while getting blog posts comments data"

		return http.StatusInternalServerError, response
	}

	response.Data = blogPosts
	response.Message = "blog posts comments data"

	return http.StatusOK, response
}
